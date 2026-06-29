// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SyncHosts is the golang structure of table sync_hosts for DAO operations like Where/Data.
type SyncHosts struct {
	g.Meta      `orm:"table:sync_hosts, do:true"`
	Id          any         //
	Name        any         //
	Hostname    any         //
	Port        any         //
	Username    any         //
	AuthMethod  any         //
	PasswordEnc any         //
	KeyId       any         //
	GroupId     any         //
	Color       any         //
	SyncUser    any         //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
