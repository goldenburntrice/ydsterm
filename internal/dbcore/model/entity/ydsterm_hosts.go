// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermHosts is the golang structure for table ydsterm_hosts.
type YdstermHosts struct {
	Id          string      `json:"id"          orm:"id"           description:""`
	Name        string      `json:"name"        orm:"name"         description:""`
	Hostname    string      `json:"hostname"    orm:"hostname"     description:""`
	Port        int         `json:"port"        orm:"port"         description:""`
	Username    string      `json:"username"    orm:"username"     description:""`
	AuthMethod  string      `json:"authMethod"  orm:"auth_method"  description:""`
	PasswordEnc string      `json:"passwordEnc" orm:"password_enc" description:""`
	KeyId       string      `json:"keyId"       orm:"key_id"       description:""`
	GroupId     string      `json:"groupId"     orm:"group_id"     description:""`
	Color       string      `json:"color"       orm:"color"        description:""`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:""`
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:""`
	SyncUser    string      `json:"syncUser"    orm:"sync_user"    description:""`
}
