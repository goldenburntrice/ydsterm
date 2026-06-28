package service

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/crypto/ssh"

	"ydsterm/internal/dbcore/dao"
	"ydsterm/internal/types"
)

// ---------- shared types ----------

type FileEntry struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	IsDir   bool   `json:"isDir"`
	Size    int64  `json:"size"`
	Mode    string `json:"mode"`
	ModTime string `json:"modTime"`
}

type SftpTransferPayload struct {
	TransferId string
	BytesDone  int64
	BytesTotal int64
	Done       bool
	Err        string
}

type sftpSession struct {
	ID        string
	HostID    string
	client    *ssh.Client
	sftpCli   *sftp.Client
	mu        sync.Mutex
	createdAt time.Time
}

// ---------- SFTPService ----------

type SFTPServiceImpl struct {
	sessions    map[string]*sftpSession
	mu          sync.Mutex
	cancelFlags map[string]chan struct{}
	cancelMu    sync.Mutex
}

func NewSFTPService() *SFTPServiceImpl {
	return &SFTPServiceImpl{
		sessions:    make(map[string]*sftpSession),
		cancelFlags: make(map[string]chan struct{}),
	}
}

func (s *SFTPServiceImpl) Connect(hostID string) (sessionID string, err error) {
	log.Printf("[sftp] Connect(%s)", hostID)
	host, err := dao.YdstermHosts.Get(nil, hostID)
	if err != nil {
		return "", fmt.Errorf("host not found: %s", hostID)
	}

	cfg, err := buildSSHConfig(host)
	if err != nil {
		return "", err
	}
	addr := fmt.Sprintf("%s:%d", host.Hostname, host.Port)
	client, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return "", fmt.Errorf("dial %s: %w", addr, err)
	}

	sftpCli, err := sftp.NewClient(client)
	if err != nil {
		client.Close()
		return "", fmt.Errorf("sftp client: %w", err)
	}

	sessionID = types.NewID()
	ss := &sftpSession{
		ID:        sessionID,
		HostID:    hostID,
		client:    client,
		sftpCli:   sftpCli,
		createdAt: time.Now(),
	}
	s.mu.Lock()
	s.sessions[sessionID] = ss
	s.mu.Unlock()
	log.Printf("[sftp] connected %s", sessionID)
	return sessionID, nil
}

func (s *SFTPServiceImpl) Disconnect(sessionID string) error {
	s.mu.Lock()
	ss, ok := s.sessions[sessionID]
	if ok {
		delete(s.sessions, sessionID)
	}
	s.mu.Unlock()
	if !ok {
		return nil
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if ss.sftpCli != nil {
		ss.sftpCli.Close()
	}
	if ss.client != nil {
		ss.client.Close()
	}
	return nil
}

func (s *SFTPServiceImpl) List(sessionID, path string) ([]FileEntry, error) {
	ss, err := s.getSession(sessionID)
	if err != nil {
		return nil, err
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if path == "" {
		path = "/"
	}

	infos, err := ss.sftpCli.ReadDir(path)
	if err != nil {
		return nil, err
	}
	sort.Slice(infos, func(i, j int) bool {
		if infos[i].IsDir() != infos[j].IsDir() {
			return infos[i].IsDir()
		}
		return strings.ToLower(infos[i].Name()) < strings.ToLower(infos[j].Name())
	})

	entries := make([]FileEntry, 0, len(infos))
	prefix := strings.TrimRight(path, "/")
	for _, fi := range infos {
		entries = append(entries, FileEntry{
			Name:    fi.Name(),
			Path:    prefix + "/" + fi.Name(),
			IsDir:   fi.IsDir(),
			Size:    fi.Size(),
			Mode:    fi.Mode().String(),
			ModTime: fi.ModTime().Format("2006-01-02 15:04"),
		})
	}
	return entries, nil
}

func (s *SFTPServiceImpl) Download(sessionID, path string) (string, error) {
	ss, err := s.getSession(sessionID)
	if err != nil {
		return "", err
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()

	f, err := ss.sftpCli.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (s *SFTPServiceImpl) Upload(sessionID, path, encoded string) error {
	ss, err := s.getSession(sessionID)
	if err != nil {
		return err
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}

	f, err := ss.sftpCli.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

func (s *SFTPServiceImpl) Mkdir(sessionID, path string) error {
	ss, err := s.getSession(sessionID)
	if err != nil {
		return err
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.sftpCli.Mkdir(path)
}

func (s *SFTPServiceImpl) Delete(sessionID, path string) error {
	ss, err := s.getSession(sessionID)
	if err != nil {
		return err
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()

	info, err := ss.sftpCli.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return s.removeDir(ss.sftpCli, path)
	}
	return ss.sftpCli.Remove(path)
}

func (s *SFTPServiceImpl) Rename(sessionID, oldPath, newPath string) error {
	ss, err := s.getSession(sessionID)
	if err != nil {
		return err
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.sftpCli.Rename(oldPath, newPath)
}

func (s *SFTPServiceImpl) HomeDir(sessionID string) (string, error) {
	ss, err := s.getSession(sessionID)
	if err != nil {
		return "", err
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.sftpCli.Getwd()
}

// ---------- chunked transfer with progress ----------

const chunkSize = 32 * 1024

func (s *SFTPServiceImpl) CancelTransfer(transferId string) {
	s.cancelMu.Lock()
	if ch, ok := s.cancelFlags[transferId]; ok {
		close(ch)
	}
	s.cancelMu.Unlock()
}

func (s *SFTPServiceImpl) registerTransfer(transferId string) chan struct{} {
	s.cancelMu.Lock()
	defer s.cancelMu.Unlock()
	ch := make(chan struct{})
	s.cancelFlags[transferId] = ch
	return ch
}

func (s *SFTPServiceImpl) unregisterTransfer(transferId string) {
	s.cancelMu.Lock()
	delete(s.cancelFlags, transferId)
	s.cancelMu.Unlock()
}

func emitProgress(transferId string, done int64, total int64, doneFlag bool, err string) {
	application.Get().Event.Emit("sftp-transfer-progress", SftpTransferPayload{
		TransferId: transferId,
		BytesDone:  done,
		BytesTotal: total,
		Done:       doneFlag,
		Err:        err,
	})
}

func (s *SFTPServiceImpl) UploadFile(sessionID, localPath, remotePath, transferId string) error {
	f, err := os.Open(localPath)
	if err != nil {
		emitProgress(transferId, 0, 0, true, err.Error())
		return err
	}
	defer f.Close()
	info, _ := f.Stat()
	total := info.Size()

	cancel := s.registerTransfer(transferId)
	defer s.unregisterTransfer(transferId)

	ss, err := s.getSession(sessionID)
	if err != nil {
		emitProgress(transferId, 0, total, true, err.Error())
		return err
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()

	rf, err := ss.sftpCli.Create(remotePath)
	if err != nil {
		emitProgress(transferId, 0, total, true, err.Error())
		return err
	}
	defer rf.Close()

	buf := make([]byte, chunkSize)
	var done int64
	for {
		select {
		case <-cancel:
			rf.Close()
			ss.sftpCli.Remove(remotePath)
			emitProgress(transferId, done, total, true, "cancelled")
			return nil
		default:
		}
		n, err := f.Read(buf)
		if n > 0 {
			if _, werr := rf.Write(buf[:n]); werr != nil {
				emitProgress(transferId, done, total, true, werr.Error())
				return werr
			}
			done += int64(n)
			emitProgress(transferId, done, total, false, "")
		}
		if err != nil {
			if err == io.EOF {
				emitProgress(transferId, total, total, true, "")
				return nil
			}
			emitProgress(transferId, done, total, true, err.Error())
			return err
		}
	}
}

func (s *SFTPServiceImpl) DownloadFile(sessionID, remotePath, localPath, transferId string) error {
	ss, err := s.getSession(sessionID)
	if err != nil {
		emitProgress(transferId, 0, 0, true, err.Error())
		return err
	}
	ss.mu.Lock()
	rf, err := ss.sftpCli.Open(remotePath)
	if err != nil {
		ss.mu.Unlock()
		emitProgress(transferId, 0, 0, true, err.Error())
		return err
	}
	defer rf.Close()
	info, _ := rf.Stat()
	total := info.Size()
	ss.mu.Unlock()

	cancel := s.registerTransfer(transferId)
	defer s.unregisterTransfer(transferId)

	f, err := os.Create(localPath)
	if err != nil {
		emitProgress(transferId, 0, total, true, err.Error())
		return err
	}
	defer f.Close()

	buf := make([]byte, chunkSize)
	var done int64
	for {
		select {
		case <-cancel:
			f.Close()
			os.Remove(localPath)
			emitProgress(transferId, done, total, true, "cancelled")
			return nil
		default:
		}
		n, err := rf.Read(buf)
		if n > 0 {
			if _, werr := f.Write(buf[:n]); werr != nil {
				emitProgress(transferId, done, total, true, werr.Error())
				return werr
			}
			done += int64(n)
			emitProgress(transferId, done, total, false, "")
		}
		if err != nil {
			if err == io.EOF {
				emitProgress(transferId, total, total, true, "")
				return nil
			}
			emitProgress(transferId, done, total, true, err.Error())
			return err
		}
	}
}

// ---------- local file methods ----------

func (s *SFTPServiceImpl) ListLocal(path string) ([]FileEntry, error) {
	if path == "" {
		path = homeDir()
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	out := make([]FileEntry, 0, len(entries))
	for _, e := range entries {
		info, _ := e.Info()
		size := int64(0)
		mode := ""
		modTime := ""
		if info != nil {
			size = info.Size()
			mode = info.Mode().String()
			modTime = info.ModTime().Format("2006-01-02 15:04")
		}
		out = append(out, FileEntry{
			Name:    e.Name(),
			Path:    filepath.Join(path, e.Name()),
			IsDir:   e.IsDir(),
			Size:    size,
			Mode:    mode,
			ModTime: modTime,
		})
	}
	return out, nil
}

func (s *SFTPServiceImpl) ReadLocal(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (s *SFTPServiceImpl) WriteLocal(path, encoded string) error {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (s *SFTPServiceImpl) MkdirLocal(path string) error {
	return os.MkdirAll(path, 0755)
}

func (s *SFTPServiceImpl) DeleteLocal(path string) error {
	return os.RemoveAll(path)
}

func (s *SFTPServiceImpl) RenameLocal(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

func (s *SFTPServiceImpl) HomeLocal() string {
	return homeDir()
}

// ---------- helpers ----------

func (s *SFTPServiceImpl) getSession(id string) (*sftpSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ss, ok := s.sessions[id]
	if !ok {
		return nil, fmt.Errorf("sftp session not found: %s", id)
	}
	return ss, nil
}

func (s *SFTPServiceImpl) removeDir(cli *sftp.Client, path string) error {
	entries, err := cli.ReadDir(path)
	if err != nil {
		return err
	}
	for _, e := range entries {
		p := strings.TrimRight(path, "/") + "/" + e.Name()
		if e.IsDir() {
			if err := s.removeDir(cli, p); err != nil {
				return err
			}
		} else {
			if err := cli.Remove(p); err != nil {
				return err
			}
		}
	}
	return cli.RemoveDirectory(path)
}

func (s *SFTPServiceImpl) ParentDir(p string, isLocal bool) string {
	if isLocal {
		return filepath.Dir(p)
	}
	if p == "" || p == "/" {
		return "/"
	}
	parent := path.Dir(p)
	if parent == "." {
		return "/"
	}
	return parent
}

func (s *SFTPServiceImpl) BaseName(p string, isLocal bool) string {
	if isLocal {
		return filepath.Base(p)
	}
	return path.Base(p)
}

func homeDir() string {
	if h, err := os.UserHomeDir(); err == nil {
		return h
	}
	return "/"
}
