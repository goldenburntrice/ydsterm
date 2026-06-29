// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermSnippets is the golang structure for table ydsterm_snippets.
type YdstermSnippets struct {
	Id        string      `json:"id"        orm:"id"         description:""`
	Name      string      `json:"name"      orm:"name"       description:""`
	Content   string      `json:"content"   orm:"content"    description:""`
	Language  string      `json:"language"  orm:"language"   description:""`
	FolderId  string      `json:"folderId"  orm:"folder_id"  description:""`
	SortOrder int         `json:"sortOrder" orm:"sort_order" description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
	SyncUser  string      `json:"syncUser"  orm:"sync_user"  description:""`
}
