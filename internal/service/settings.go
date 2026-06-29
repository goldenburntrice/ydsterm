package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"ydsterm/internal/dbcore/dao"
	"ydsterm/internal/types"
)

type SettingsServiceImpl struct{}

func NewSettingsService() *SettingsServiceImpl { return &SettingsServiceImpl{} }

func (s *SettingsServiceImpl) Get(key string) (string, error) {
	return dao.YdstermSettings.Get(context.Background(), key)
}

func (s *SettingsServiceImpl) Set(key, value string) error {
	return dao.YdstermSettings.Set(context.Background(), key, value)
}

func (s *SettingsServiceImpl) VerifyUser(serverAddr, serverKey, username, password string) (*types.VerifyResponse, error) {
	reqBody := types.VerifyRequest{Username: username, Password: password}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := serverAddr + "/api/auth/verify"
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Server-Key", serverKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("无法连接到服务器: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("密码错误")
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("服务器错误: %s", string(respBody))
	}

	var result types.VerifyResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}
	result.Created = resp.StatusCode == 201
	result.PasswordHash = hashPassword(password)
	return &result, nil
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}
