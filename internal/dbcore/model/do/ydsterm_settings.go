// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermSettings is the golang structure of table ydsterm_settings for DAO operations like Where/Data.
type YdstermSettings struct {
	g.Meta    `orm:"table:ydsterm_settings, do:true"`
	Key       any         //
	Value     any         //
	UpdatedAt *gtime.Time //
}
