package types

import "time"

type SyncUser struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	IsActive   bool      `json:"isActive"`
	LastSyncAt *time.Time `json:"lastSyncAt,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type VerifyRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VerifyResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Created  bool   `json:"created"`
}

type SyncData struct {
	Hosts           SyncTableData `json:"hosts"`
	HostGroups      SyncTableData `json:"host_groups"`
	Keys            SyncTableData `json:"keys"`
	Snippets        SyncTableData `json:"snippets"`
	SnippetFolders  SyncTableData `json:"snippet_folders"`
	PortForwards    SyncTableData `json:"port_forwards"`
}

type SyncTableData struct {
	Upsert []SyncRecord `json:"upsert"`
	Delete []string     `json:"delete"`
}

type SyncRecord struct {
	ID        string                 `json:"id"`
	Data      map[string]interface{} `json:"data"`
	UpdatedAt string                 `json:"updated_at"`
}

type SyncPushRequest struct {
	LastSyncAt string   `json:"last_sync_at"`
	Data       SyncData `json:"data"`
}

type SyncPullResponse struct {
	ServerTime string   `json:"server_time"`
	Data       SyncData `json:"data"`
}
