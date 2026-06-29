package types

import "time"

type SyncUser struct {
	ID         string     `json:"id"`
	Username   string     `json:"username"`
	IsActive   bool       `json:"isActive"`
	LastSyncAt *time.Time `json:"lastSyncAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

type VerifyRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VerifyResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Created      bool   `json:"created"`
	PasswordHash string `json:"passwordHash"`
}
