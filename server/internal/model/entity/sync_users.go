// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SyncUsers is the golang structure for table sync_users.
type SyncUsers struct {
	Id           string      `json:"id"           orm:"id"            ` //
	Username     string      `json:"username"     orm:"username"      ` //
	PasswordHash string      `json:"passwordHash" orm:"password_hash" ` //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` //
}
