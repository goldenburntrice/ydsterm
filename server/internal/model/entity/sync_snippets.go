// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SyncSnippets is the golang structure for table sync_snippets.
type SyncSnippets struct {
	Id        string      `json:"id"        orm:"id"         ` //
	Name      string      `json:"name"      orm:"name"       ` //
	Content   string      `json:"content"   orm:"content"    ` //
	Language  string      `json:"language"  orm:"language"   ` //
	FolderId  string      `json:"folderId"  orm:"folder_id"  ` //
	SortOrder int         `json:"sortOrder" orm:"sort_order" ` //
	SyncUser  string      `json:"syncUser"  orm:"sync_user"  ` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` //
}
