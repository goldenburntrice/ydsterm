// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermSnippetFolders is the golang structure for table ydsterm_snippet_folders.
type YdstermSnippetFolders struct {
	Id        string      `json:"id"        orm:"id"         description:""`
	Name      string      `json:"name"      orm:"name"       description:""`
	ParentId  string      `json:"parentId"  orm:"parent_id"  description:""`
	SortOrder int         `json:"sortOrder" orm:"sort_order" description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
