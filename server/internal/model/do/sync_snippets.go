// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SyncSnippets is the golang structure of table sync_snippets for DAO operations like Where/Data.
type SyncSnippets struct {
	g.Meta    `orm:"table:sync_snippets, do:true"`
	Id        any         //
	Name      any         //
	Content   any         //
	Language  any         //
	FolderId  any         //
	SortOrder any         //
	SyncUser  any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
