package controller

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"golang.org/x/crypto/bcrypt"

	"ydsterm-server/internal/dao"
)

type AuthController struct{}

var Auth = AuthController{}

type VerifyRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VerifyResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (c *AuthController) Verify(r *ghttp.Request) {
	ctx := r.Context()
	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": "invalid request body"})
		return
	}
	if req.Username == "" || req.Password == "" {
		r.Response.WriteJsonExit(g.Map{"error": "username and password required"})
		return
	}

	record, err := dao.SyncUsers.Ctx(ctx).Where("username", req.Username).One()
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"error": "database error"})
		return
	}

	if record.IsEmpty() {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			r.Response.WriteJsonExit(g.Map{"error": "failed to hash password"})
			return
		}
		b := make([]byte, 16)
		rand.Read(b)
		id := hex.EncodeToString(b)
		_, err = dao.SyncUsers.Ctx(ctx).Data(g.Map{
			"id":            id,
			"username":      req.Username,
			"password_hash": string(hash),
		}).Insert()
		if err != nil {
			r.Response.WriteJsonExit(g.Map{"error": "failed to create user"})
			return
		}
		r.Response.WriteJson(VerifyResponse{ID: id, Username: req.Username})
		return
	}

	hash := record["password_hash"].String()
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		r.Response.WriteStatusExit(401, g.Map{"error": "invalid password"})
		return
	}

	r.Response.WriteJson(VerifyResponse{
		ID:       record["id"].String(),
		Username: record["username"].String(),
	})
}

type SyncController struct{}

var Sync = SyncController{}

type SyncRecord struct {
	ID        string                 `json:"id"`
	Data      map[string]interface{} `json:"data"`
	UpdatedAt string                 `json:"updated_at"`
}

type SyncTableData struct {
	Upsert []SyncRecord `json:"upsert"`
	Delete []string     `json:"delete"`
}

type SyncData struct {
	Hosts          SyncTableData `json:"hosts"`
	HostGroups     SyncTableData `json:"host_groups"`
	Keys           SyncTableData `json:"keys"`
	Snippets       SyncTableData `json:"snippets"`
	SnippetFolders SyncTableData `json:"snippet_folders"`
	PortForwards   SyncTableData `json:"port_forwards"`
}

type PullResponse struct {
	ServerTime string   `json:"server_time"`
	Data       SyncData `json:"data"`
}

type PushRequest struct {
	LastSyncAt string   `json:"last_sync_at"`
	Data       SyncData `json:"data"`
}

func (c *SyncController) Pull(r *ghttp.Request) {
	ctx := r.Context()
	username := r.GetQuery("user").String()
	if username == "" {
		r.Response.WriteJsonExit(g.Map{"error": "user parameter required"})
		return
	}
	since := r.GetQuery("since").String()

	serverTime := time.Now().UTC().Format(time.RFC3339)
	data := SyncData{}

	tables := []struct {
		table string
		out   *SyncTableData
	}{
		{"sync_hosts", &data.Hosts},
		{"sync_host_groups", &data.HostGroups},
		{"sync_keys", &data.Keys},
		{"sync_snippets", &data.Snippets},
		{"sync_snippet_folders", &data.SnippetFolders},
		{"sync_port_forwards", &data.PortForwards},
	}

	for _, t := range tables {
		records, err := queryTable(ctx, t.table, username, since)
		if err != nil {
			r.Response.WriteJsonExit(g.Map{"error": "database error"})
			return
		}
		t.out.Upsert = records
		t.out.Delete = []string{}
	}

	r.Response.WriteJson(PullResponse{ServerTime: serverTime, Data: data})
}

func (c *SyncController) Push(r *ghttp.Request) {
	ctx := r.Context()
	username := r.GetQuery("user").String()
	if username == "" {
		r.Response.WriteJsonExit(g.Map{"error": "user parameter required"})
		return
	}

	var req PushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": "invalid request body"})
		return
	}

	tables := []struct {
		table string
		data  SyncTableData
	}{
		{"sync_hosts", req.Data.Hosts},
		{"sync_host_groups", req.Data.HostGroups},
		{"sync_keys", req.Data.Keys},
		{"sync_snippets", req.Data.Snippets},
		{"sync_snippet_folders", req.Data.SnippetFolders},
		{"sync_port_forwards", req.Data.PortForwards},
	}

	for _, t := range tables {
		if err := upsertTable(ctx, t.table, username, t.data); err != nil {
			r.Response.WriteJsonExit(g.Map{"error": "database error"})
			return
		}
	}

	r.Response.WriteJson(g.Map{"status": "ok"})
}

func queryTable(ctx context.Context, table, username, since string) ([]SyncRecord, error) {
	m := g.DB().Model(table).Ctx(ctx).Where("sync_user", username)
	if since != "" {
		m = m.WhereGT("updated_at", since)
	}

	records, err := m.All()
	if err != nil {
		return nil, err
	}

	var result []SyncRecord
	for _, r := range records {
		data := r.Map()
		id := ""
		if v, ok := data["id"]; ok && v != nil {
			id = v.(string)
		}
		updatedAt := ""
		if t, ok := data["updated_at"]; ok && t != nil {
			updatedAt = t.(string)
		}
		result = append(result, SyncRecord{
			ID:        id,
			Data:      data,
			UpdatedAt: updatedAt,
		})
	}
	if result == nil {
		result = []SyncRecord{}
	}
	return result, nil
}

func upsertTable(ctx context.Context, table, username string, data SyncTableData) error {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	for _, record := range data.Upsert {
		record.Data["sync_user"] = username
		record.Data["updated_at"] = now
		if _, ok := record.Data["created_at"]; !ok {
			record.Data["created_at"] = now
		}

		_, err := g.DB().Model(table).Ctx(ctx).Data(record.Data).Save()
		if err != nil {
			return err
		}
	}

	for _, id := range data.Delete {
		_, err := g.DB().Model(table).Ctx(ctx).Where("id", id).Where("sync_user", username).Delete()
		if err != nil {
			return err
		}
	}

	return nil
}
