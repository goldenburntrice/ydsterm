// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermSnippets is the golang structure of table ydsterm_snippets for DAO operations like Where/Data.
type YdstermSnippets struct {
	g.Meta    `orm:"table:ydsterm_snippets, do:true"`
	Id        any         //
	Name      any         //
	Content   any         //
	Language  any         //
	FolderId  any         //
	SortOrder any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
