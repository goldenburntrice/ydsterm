// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermKeys is the golang structure of table ydsterm_keys for DAO operations like Where/Data.
type YdstermKeys struct {
	g.Meta        `orm:"table:ydsterm_keys, do:true"`
	Id            any         //
	Name          any         //
	PrivateKeyEnc any         //
	PublicKey     any         //
	PassphraseEnc any         //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
	SyncUser      any         //
}
