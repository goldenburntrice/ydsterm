// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermSchemaMigrations is the golang structure for table ydsterm_schema_migrations.
type YdstermSchemaMigrations struct {
	Version   int         `json:"version"   orm:"version"    description:""`
	AppliedAt *gtime.Time `json:"appliedAt" orm:"applied_at" description:""`
}
