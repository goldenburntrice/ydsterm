// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermHosts is the golang structure of table ydsterm_hosts for DAO operations like Where/Data.
type YdstermHosts struct {
	g.Meta      `orm:"table:ydsterm_hosts, do:true"`
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
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
	SyncUser    any         //
}
