package types

import "time"

type Host struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Hostname    string     `json:"hostname"`
	Port        int        `json:"port"`
	Username    string     `json:"username"`
	AuthMethod  string     `json:"authMethod"`
	PasswordEnc string     `json:"-"`
	KeyID       string     `json:"keyId"`
	GroupID     string     `json:"groupId"`
	Color       string     `json:"color"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type HostCreateInput struct {
	Name       string `json:"name"`
	Hostname   string `json:"hostname"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	AuthMethod string `json:"authMethod"`
	Password   string `json:"password,omitempty"`
	KeyID      string `json:"keyId"`
	GroupID    string `json:"groupId"`
	Color      string `json:"color"`
}

type HostUpdateInput struct {
	ID         string  `json:"id"`
	Name       *string `json:"name,omitempty"`
	Hostname   *string `json:"hostname,omitempty"`
	Port       *int    `json:"port,omitempty"`
	Username   *string `json:"username,omitempty"`
	AuthMethod *string `json:"authMethod,omitempty"`
	Password   *string `json:"password,omitempty"`
	KeyID      *string `json:"keyId,omitempty"`
	GroupID    *string `json:"groupId,omitempty"`
	Color      *string `json:"color,omitempty"`
}

type Key struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	PrivateKeyEnc string    `json:"-"`
	PublicKey     string    `json:"publicKey"`
	PassphraseEnc string    `json:"-"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type KeyCreateInput struct {
	Name       string `json:"name"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey,omitempty"`
	Passphrase string `json:"passphrase,omitempty"`
}

type KeyUpdateInput struct {
	ID         string  `json:"id"`
	Name       *string `json:"name,omitempty"`
	PrivateKey *string `json:"privateKey,omitempty"`
	PublicKey  *string `json:"publicKey,omitempty"`
	Passphrase *string `json:"passphrase,omitempty"`
}

type HostGroup struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ParentID  string    `json:"parentId"`
	SortOrder int       `json:"sortOrder"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type HostGroupCreateInput struct {
	Name     string `json:"name"`
	ParentID string `json:"parentId"`
}

type HostGroupUpdateInput struct {
	ID       string  `json:"id"`
	Name     *string `json:"name,omitempty"`
	ParentID *string `json:"parentId,omitempty"`
}
