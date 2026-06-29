package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/wailsapp/wails/v3/pkg/application"

	"ydsterm/internal/crypto"
	"ydsterm/internal/dbcore"
	"ydsterm/internal/dbcore/dao"
)

type SyncServiceImpl struct{}

func NewSyncService() *SyncServiceImpl { return &SyncServiceImpl{} }

type syncModule struct {
	name       string
	table      string
	sensFields []string
}

var syncModules = []syncModule{
	{"host_groups", "ydsterm_host_groups", nil},
	{"snippet_folders", "ydsterm_snippet_folders", nil},
	{"hosts", "ydsterm_hosts", []string{"password_enc"}},
	{"snippets", "ydsterm_snippets", nil},
	{"keys", "ydsterm_keys", []string{"private_key_enc", "passphrase_enc"}},
	{"port_forwards", "ydsterm_port_forwards", nil},
}

type SyncProgressPayload struct {
	Module string `json:"module"`
	Done   bool   `json:"done"`
	Error  string `json:"error"`
}

func init() {
	application.RegisterEvent[SyncProgressPayload]("sync-progress")
}

func (s *SyncServiceImpl) emitProgress(module string, done bool, errStr string) {
	application.Get().Event.Emit("sync-progress", SyncProgressPayload{
		Module: module,
		Done:   done,
		Error:  errStr,
	})
}

func (s *SyncServiceImpl) getServerConfig() (string, string, error) {
	addr, _ := dao.YdstermSettings.Get(context.Background(), "server_addr")
	key, _ := dao.YdstermSettings.Get(context.Background(), "server_key")
	if addr == "" || key == "" {
		return "", "", fmt.Errorf("server not configured")
	}
	return addr, key, nil
}

func (s *SyncServiceImpl) Push(module string, data map[string]interface{}) error {
	if dbcore.GetActiveUser() == "LOCALUSER" {
		return nil
	}

	addr, serverKey, err := s.getServerConfig()
	if err != nil {
		return err
	}

	syncData := s.decryptForPush(module, data)

	jsonBytes, err := json.Marshal(syncData)
	if err != nil {
		return err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	md5Sign := crypto.ComputeMD5Sign(syncData, timestamp, serverKey)

	encrypted, err := crypto.EncryptWithKey(serverKey, string(jsonBytes))
	if err != nil {
		return err
	}

	user := dbcore.GetActiveUser()
	url := fmt.Sprintf("%s/sync2cloud/%s?user=%s", addr, module, user)

	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(encrypted)))
	if err != nil {
		return err
	}
	req.Header.Set("X-MD5", md5Sign)
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("push failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("push failed: status %d", resp.StatusCode)
	}

	return nil
}

func (s *SyncServiceImpl) Delete(module string, id string) error {
	if dbcore.GetActiveUser() == "LOCALUSER" {
		return nil
	}

	addr, serverKey, err := s.getServerConfig()
	if err != nil {
		return err
	}

	data := map[string]interface{}{"id": id}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	md5Sign := crypto.ComputeMD5Sign(data, timestamp, serverKey)

	encrypted, err := crypto.EncryptWithKey(serverKey, string(jsonBytes))
	if err != nil {
		return err
	}

	user := dbcore.GetActiveUser()
	url := fmt.Sprintf("%s/sync2cloud/%s?user=%s", addr, module, user)

	req, err := http.NewRequest("DELETE", url, bytes.NewReader([]byte(encrypted)))
	if err != nil {
		return err
	}
	req.Header.Set("X-MD5", md5Sign)
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("delete failed: status %d", resp.StatusCode)
	}

	return nil
}

func (s *SyncServiceImpl) PullAll() error {
	if dbcore.GetActiveUser() == "LOCALUSER" {
		return nil
	}

	addr, serverKey, err := s.getServerConfig()
	if err != nil {
		return err
	}

	user := dbcore.GetActiveUser()
	client := &http.Client{Timeout: 30 * time.Second}

	for _, mod := range syncModules {
		s.emitProgress(mod.name, false, "")

		url := fmt.Sprintf("%s/pull2cli/%s?user=%s&server_key=%s", addr, mod.name, user, serverKey)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			errMsg := fmt.Sprintf("pull %s failed: %v", mod.name, err)
			s.emitProgress(mod.name, false, errMsg)
			return fmt.Errorf(errMsg)
		}

		resp, err := client.Do(req)
		if err != nil {
			errMsg := fmt.Sprintf("pull %s failed: %v", mod.name, err)
			s.emitProgress(mod.name, false, errMsg)
			return fmt.Errorf(errMsg)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			errMsg := fmt.Sprintf("pull %s read failed: %v", mod.name, err)
			s.emitProgress(mod.name, false, errMsg)
			return fmt.Errorf(errMsg)
		}

		decrypted, err := crypto.DecryptWithKey(serverKey, string(body))
		if err != nil {
			errMsg := fmt.Sprintf("decrypt %s failed: %v", mod.name, err)
			s.emitProgress(mod.name, false, errMsg)
			return fmt.Errorf(errMsg)
		}

		var records []map[string]interface{}
		if err := json.Unmarshal([]byte(decrypted), &records); err != nil {
			errMsg := fmt.Sprintf("decode %s failed: %v", mod.name, err)
			s.emitProgress(mod.name, false, errMsg)
			return fmt.Errorf(errMsg)
		}

		if err := s.applyPullData(mod, records, serverKey); err != nil {
			errMsg := fmt.Sprintf("apply %s failed: %v", mod.name, err)
			s.emitProgress(mod.name, false, errMsg)
			return fmt.Errorf(errMsg)
		}

		s.emitProgress(mod.name, true, "")
	}

	return nil
}

func (s *SyncServiceImpl) applyPullData(mod syncModule, records []map[string]interface{}, serverKey string) error {
	ctx := dbcore.WithActiveUser(context.Background())

	s.encryptForLocal(mod, records, serverKey)

	serverIDs := make(map[string]bool)
	for _, r := range records {
		id, _ := r["id"].(string)
		if id == "" {
			continue
		}
		serverIDs[id] = true
		delete(r, "sync_user")

		r["updated_at"] = time.Now().Format("2006-01-02 15:04:05")
		if _, err := g.DB().Model(mod.table).Ctx(ctx).Data(r).Save(); err != nil {
			return err
		}
	}

	allRecords, err := g.DB().Model(mod.table).Ctx(ctx).All()
	if err != nil {
		return err
	}

	for _, record := range allRecords {
		id := record["id"].String()
		if !serverIDs[id] {
			deletedAt := record["deleted_at"]
			if deletedAt == nil || deletedAt.IsEmpty() {
				if _, err := g.DB().Model(mod.table).Ctx(ctx).Where("id", id).Update(g.Map{
					"deleted_at": time.Now().Format("2006-01-02 15:04:05"),
				}); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *SyncServiceImpl) decryptForPush(module string, data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		result[k] = v
	}

	for _, mod := range syncModules {
		if mod.name != module {
			continue
		}
		for _, field := range mod.sensFields {
			if enc, ok := result[field].(string); ok && enc != "" {
				plain, err := crypto.Decrypt(enc)
				if err == nil {
					result[field] = plain
				}
			}
		}
	}

	return result
}

func (s *SyncServiceImpl) encryptForLocal(mod syncModule, records []map[string]interface{}, serverKey string) {
	if len(mod.sensFields) == 0 {
		return
	}

	for _, r := range records {
		for _, field := range mod.sensFields {
			if plain, ok := r[field].(string); ok && plain != "" {
				enc, err := crypto.EncryptWithKey(serverKey, plain)
				if err == nil {
					r[field] = enc
				}
			}
		}
	}
}

func (s *SyncServiceImpl) PushHost(data map[string]interface{}) {
	if err := s.Push("hosts", data); err != nil {
		log.Printf("sync push hosts: %v", err)
	}
}

func (s *SyncServiceImpl) PushHostGroup(data map[string]interface{}) {
	if err := s.Push("host_groups", data); err != nil {
		log.Printf("sync push host_groups: %v", err)
	}
}

func (s *SyncServiceImpl) PushKey(data map[string]interface{}) {
	if err := s.Push("keys", data); err != nil {
		log.Printf("sync push keys: %v", err)
	}
}

func (s *SyncServiceImpl) PushSnippet(data map[string]interface{}) {
	if err := s.Push("snippets", data); err != nil {
		log.Printf("sync push snippets: %v", err)
	}
}

func (s *SyncServiceImpl) PushSnippetFolder(data map[string]interface{}) {
	if err := s.Push("snippet_folders", data); err != nil {
		log.Printf("sync push snippet_folders: %v", err)
	}
}

func (s *SyncServiceImpl) PushPortForward(data map[string]interface{}) {
	if err := s.Push("port_forwards", data); err != nil {
		log.Printf("sync push port_forwards: %v", err)
	}
}

func (s *SyncServiceImpl) DeleteHost(id string) {
	if err := s.Delete("hosts", id); err != nil {
		log.Printf("sync delete hosts: %v", err)
	}
}

func (s *SyncServiceImpl) DeleteHostGroup(id string) {
	if err := s.Delete("host_groups", id); err != nil {
		log.Printf("sync delete host_groups: %v", err)
	}
}

func (s *SyncServiceImpl) DeleteKey(id string) {
	if err := s.Delete("keys", id); err != nil {
		log.Printf("sync delete keys: %v", err)
	}
}

func (s *SyncServiceImpl) DeleteSnippet(id string) {
	if err := s.Delete("snippets", id); err != nil {
		log.Printf("sync delete snippets: %v", err)
	}
}

func (s *SyncServiceImpl) DeleteSnippetFolder(id string) {
	if err := s.Delete("snippet_folders", id); err != nil {
		log.Printf("sync delete snippet_folders: %v", err)
	}
}

func (s *SyncServiceImpl) DeletePortForward(id string) {
	if err := s.Delete("port_forwards", id); err != nil {
		log.Printf("sync delete port_forwards: %v", err)
	}
}
