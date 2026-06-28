package service

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"

	"ydsterm/internal/crypto"
	"ydsterm/internal/dbcore/dao"
	"ydsterm/internal/dbcore/model/entity"
)

func buildSSHConfig(host *entity.YdstermHosts) (*ssh.ClientConfig, error) {
	cfg := &ssh.ClientConfig{
		User:            host.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
	}

	switch host.AuthMethod {
	case "password":
		pw, err := crypto.Decrypt(host.PasswordEnc)
		if err != nil {
			return nil, fmt.Errorf("decrypt password: %w", err)
		}
		cfg.Auth = []ssh.AuthMethod{ssh.Password(pw)}

	case "private_key":
		if host.KeyId == "" {
			return nil, fmt.Errorf("no key selected")
		}
		key, err := dao.YdstermKeys.Get(nil, host.KeyId)
		if err != nil {
			return nil, fmt.Errorf("key not found: %w", err)
		}
		pem, err := crypto.Decrypt(key.PrivateKeyEnc)
		if err != nil {
			return nil, fmt.Errorf("decrypt key: %w", err)
		}
		var signer ssh.Signer
		if key.PassphraseEnc != "" {
			pass, _ := crypto.Decrypt(key.PassphraseEnc)
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(pem), []byte(pass))
		} else {
			signer, err = ssh.ParsePrivateKey([]byte(pem))
		}
		if err != nil {
			return nil, fmt.Errorf("parse key: %w", err)
		}
		cfg.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}

	default:
		return nil, fmt.Errorf("unknown auth method: %s", host.AuthMethod)
	}
	return cfg, nil
}
