package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func EncryptWithKey(keyHex string, plaintext string) (string, error) {
	key := []byte(keyHex)
	if len(key) != 32 {
		if len(key) > 32 {
			key = key[:32]
		} else {
			padded := make([]byte, 32)
			copy(padded, key)
			key = padded
		}
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptWithKey(keyHex string, encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}

	key := []byte(keyHex)
	if len(key) != 32 {
		if len(key) > 32 {
			key = key[:32]
		} else {
			padded := make([]byte, 32)
			copy(padded, key)
			key = padded
		}
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
