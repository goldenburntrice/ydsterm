// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SyncUsers is the golang structure of table sync_users for DAO operations like Where/Data.
type SyncUsers struct {
	g.Meta       `orm:"table:sync_users, do:true"`
	Id           any         //
	Username     any         //
	PasswordHash any         //
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
