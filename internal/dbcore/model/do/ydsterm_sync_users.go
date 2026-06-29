// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermSyncUsers is the golang structure of table ydsterm_sync_users for DAO operations like Where/Data.
type YdstermSyncUsers struct {
	g.Meta     `orm:"table:ydsterm_sync_users, do:true"`
	Id         any         //
	Username   any         //
	IsActive   any         //
	LastSyncAt *gtime.Time //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
}
