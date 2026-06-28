package service

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/crypto/ssh"

	"ydsterm/internal/dbcore/dao"
)

// ---------- event payload types (registered in main.go) ----------

// Wails v3 events: field names MUST match the JavaScript property names exactly.
// JSON tags may be ignored by Wails internal serialization.
type TermOutputPayload struct {
	SessionId string
	Data      string
	IsStderr  bool
}

type TermDisconnectedPayload struct {
	SessionId string
}

// ---------- internal session ----------

type sshSession struct {
	ID        string
	HostID    string
	HostName  string
	client    *ssh.Client
	sess      *ssh.Session
	stdin     io.WriteCloser
	mu        sync.Mutex
	createdAt time.Time
}

// ---------- TerminalService (Wails-bound) ----------

type TerminalServiceImpl struct {
	sessions map[string]*sshSession
	mu       sync.Mutex
}

func NewTerminalService() *TerminalServiceImpl {
	return &TerminalServiceImpl{sessions: make(map[string]*sshSession)}
}

// Connect opens an SSH connection + PTY, returns a session ID.
func (s *TerminalServiceImpl) Connect(hostID string) (sessionID string, err error) {
	log.Printf("[term] Connect(%s) called", hostID)

	host, err := dao.YdstermHosts.Get(nil, hostID)
	if err != nil {
		log.Printf("[term] host lookup failed: %v", err)
		return "", fmt.Errorf("host not found: %s", hostID)
	}
	log.Printf("[term] host: %s (%s@%s:%d) auth=%s",
		host.Name, host.Username, host.Hostname, host.Port, host.AuthMethod)

	sessionID = hostID + "_" + fmt.Sprintf("%d", time.Now().UnixMilli())
	sess := &sshSession{
		ID:        sessionID,
		HostID:    hostID,
		HostName:  host.Name,
		createdAt: time.Now(),
	}

	cfg, err := buildSSHConfig(host)
	if err != nil {
		log.Printf("[term] ssh config: %v", err)
		return "", err
	}

	addr := fmt.Sprintf("%s:%d", host.Hostname, host.Port)
	log.Printf("[term] dialing %s...", addr)
	client, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		log.Printf("[term] dial failed: %v", err)
		return "", fmt.Errorf("connect %s: %w", addr, err)
	}
	sess.client = client
	log.Printf("[term] connected to %s", addr)

	sshSess, err := client.NewSession()
	if err != nil {
		client.Close()
		log.Printf("[term] new session: %v", err)
		return "", fmt.Errorf("create session: %w", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := sshSess.RequestPty("xterm-256color", 120, 40, modes); err != nil {
		sshSess.Close(); client.Close()
		log.Printf("[term] request pty: %v", err)
		return "", fmt.Errorf("request pty: %w", err)
	}

	stdinPipe, err := sshSess.StdinPipe()
	if err != nil {
		sshSess.Close(); client.Close()
		log.Printf("[term] stdin pipe: %v", err)
		return "", fmt.Errorf("stdin pipe: %w", err)
	}
	sess.stdin = stdinPipe

	stdoutPipe, err := sshSess.StdoutPipe()
	if err != nil {
		sshSess.Close(); client.Close()
		log.Printf("[term] stdout pipe: %v", err)
		return "", fmt.Errorf("stdout pipe: %w", err)
	}

	stderrPipe, err := sshSess.StderrPipe()
	if err != nil {
		sshSess.Close(); client.Close()
		log.Printf("[term] stderr pipe: %v", err)
		return "", fmt.Errorf("stderr pipe: %w", err)
	}

	if err := sshSess.Shell(); err != nil {
		sshSess.Close(); client.Close()
		log.Printf("[term] start shell: %v", err)
		return "", fmt.Errorf("start shell: %w", err)
	}
	sess.sess = sshSess

	s.mu.Lock()
	s.sessions[sessionID] = sess
	s.mu.Unlock()

	go s.readPipe(sessionID, stdoutPipe, false)
	go s.readPipe(sessionID, stderrPipe, true)

	log.Printf("[term] session %s started", sessionID)
	return sessionID, nil
}

// Write sends user input to the remote PTY.
func (s *TerminalServiceImpl) Write(sessionID, data string) error {
	s.mu.Lock()
	sess, ok := s.sessions[sessionID]
	s.mu.Unlock()
	if !ok {
		return fmt.Errorf("session not found: %s", sessionID)
	}
	sess.mu.Lock()
	defer sess.mu.Unlock()
	_, err := sess.stdin.Write([]byte(data))
	return err
}

// Resize adjusts the PTY window.
func (s *TerminalServiceImpl) Resize(sessionID string, cols, rows int) error {
	s.mu.Lock()
	sess, ok := s.sessions[sessionID]
	s.mu.Unlock()
	if !ok {
		return fmt.Errorf("session not found: %s", sessionID)
	}
	sess.mu.Lock()
	defer sess.mu.Unlock()
	if sess.sess != nil {
		return sess.sess.WindowChange(rows, cols)
	}
	return nil
}

// Disconnect tears down the SSH session.
func (s *TerminalServiceImpl) Disconnect(sessionID string) error {
	log.Printf("[term] Disconnect(%s)", sessionID)
	s.mu.Lock()
	sess, ok := s.sessions[sessionID]
	if ok {
		delete(s.sessions, sessionID)
	}
	s.mu.Unlock()
	if !ok {
		return nil
	}
	sess.mu.Lock()
	defer sess.mu.Unlock()
	if sess.sess != nil {
		sess.sess.Close()
	}
	if sess.client != nil {
		sess.client.Close()
	}
	return nil
}

// ListSessions returns metadata for active sessions.
func (s *TerminalServiceImpl) ListSessions() []map[string]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]map[string]interface{}, 0, len(s.sessions))
	for _, ss := range s.sessions {
		out = append(out, map[string]interface{}{
			"id":        ss.ID,
			"hostId":    ss.HostID,
			"hostName":  ss.HostName,
			"createdAt": ss.createdAt,
		})
	}
	return out
}

// ---------- internal helpers ----------

func (s *TerminalServiceImpl) readPipe(sessionID string, r io.Reader, isStderr bool) {
	buf := make([]byte, 8192)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			payload := TermOutputPayload{
				SessionId: sessionID,
				Data:      base64.StdEncoding.EncodeToString(buf[:n]),
				IsStderr:  isStderr,
			}
			application.Get().Event.Emit("term-output", payload)
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("[term] read error [%s] stderr=%v: %v", sessionID, isStderr, err)
			}
			application.Get().Event.Emit("term-disconnected", TermDisconnectedPayload{SessionId: sessionID})
			s.Disconnect(sessionID)
			return
		}
	}
}
