// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermSchemaMigrations is the golang structure of table ydsterm_schema_migrations for DAO operations like Where/Data.
type YdstermSchemaMigrations struct {
	g.Meta    `orm:"table:ydsterm_schema_migrations, do:true"`
	Version   any         //
	AppliedAt *gtime.Time //
}
