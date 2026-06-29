package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func deriveKey() []byte {
	hostname, _ := os.Hostname()
	seed := "ydsterm-encryption-key-v1-" + hostname
	hash := sha256.Sum256([]byte(seed))
	return hash[:]
}

func Encrypt(plaintext string) (string, error) {
	key := deriveKey()
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

func Decrypt(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}

	key := deriveKey()
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

func EncryptWithKey(keyHex string, plaintext string) (string, error) {
	key := normalizeKey(keyHex)
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
	key := normalizeKey(keyHex)
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

func ComputeMD5Sign(dataMap map[string]interface{}, timestamp, serverKey string) string {
	keys := make([]string, 0, len(dataMap))
	for k := range dataMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	values := make([]string, 0, len(keys))
	for _, k := range keys {
		values = append(values, fmt.Sprintf("%v", dataMap[k]))
	}

	suffix := serverKey
	if len(suffix) > 8 {
		suffix = suffix[len(suffix)-8:]
	}

	raw := strings.Join(values, "|") + timestamp + suffix
	hash := md5.Sum([]byte(raw))
	return hex.EncodeToString(hash[:])
}

func normalizeKey(keyHex string) []byte {
	key := []byte(keyHex)
	if len(key) == 32 {
		return key
	}
	if len(key) > 32 {
		return key[:32]
	}
	padded := make([]byte, 32)
	copy(padded, key)
	return padded
}
