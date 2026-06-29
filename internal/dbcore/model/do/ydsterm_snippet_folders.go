// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermSnippetFolders is the golang structure of table ydsterm_snippet_folders for DAO operations like Where/Data.
type YdstermSnippetFolders struct {
	g.Meta    `orm:"table:ydsterm_snippet_folders, do:true"`
	Id        any         //
	Name      any         //
	ParentId  any         //
	SortOrder any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
	SyncUser  any         //
}
