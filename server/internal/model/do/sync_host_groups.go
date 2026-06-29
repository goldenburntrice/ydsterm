// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SyncHostGroups is the golang structure of table sync_host_groups for DAO operations like Where/Data.
type SyncHostGroups struct {
	g.Meta    `orm:"table:sync_host_groups, do:true"`
	Id        any         //
	Name      any         //
	ParentId  any         //
	SortOrder any         //
	SyncUser  any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
