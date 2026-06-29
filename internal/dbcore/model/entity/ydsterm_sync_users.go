// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermSyncUsers is the golang structure for table ydsterm_sync_users.
type YdstermSyncUsers struct {
	Id         string      `json:"id"         orm:"id"           description:""`
	Username   string      `json:"username"   orm:"username"     description:""`
	IsActive   int         `json:"isActive"   orm:"is_active"    description:""`
	LastSyncAt *gtime.Time `json:"lastSyncAt" orm:"last_sync_at" description:""`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"   description:""`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"   description:""`
}
