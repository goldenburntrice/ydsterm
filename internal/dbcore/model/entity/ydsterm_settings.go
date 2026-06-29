// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermSettings is the golang structure for table ydsterm_settings.
type YdstermSettings struct {
	Key       string      `json:"key"       orm:"key"        description:""`
	Value     string      `json:"value"     orm:"value"      description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
