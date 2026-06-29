// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SyncKeys is the golang structure of table sync_keys for DAO operations like Where/Data.
type SyncKeys struct {
	g.Meta        `orm:"table:sync_keys, do:true"`
	Id            any         //
	Name          any         //
	PrivateKeyEnc any         //
	PublicKey     any         //
	PassphraseEnc any         //
	SyncUser      any         //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
