// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SyncHosts is the golang structure for table sync_hosts.
type SyncHosts struct {
	Id          string      `json:"id"          orm:"id"           ` //
	Name        string      `json:"name"        orm:"name"         ` //
	Hostname    string      `json:"hostname"    orm:"hostname"     ` //
	Port        int         `json:"port"        orm:"port"         ` //
	Username    string      `json:"username"    orm:"username"     ` //
	AuthMethod  string      `json:"authMethod"  orm:"auth_method"  ` //
	PasswordEnc string      `json:"passwordEnc" orm:"password_enc" ` //
	KeyId       string      `json:"keyId"       orm:"key_id"       ` //
	GroupId     string      `json:"groupId"     orm:"group_id"     ` //
	Color       string      `json:"color"       orm:"color"        ` //
	SyncUser    string      `json:"syncUser"    orm:"sync_user"    ` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   ` //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   ` //
}
