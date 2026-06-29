// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SyncPortForwards is the golang structure of table sync_port_forwards for DAO operations like Where/Data.
type SyncPortForwards struct {
	g.Meta       `orm:"table:sync_port_forwards, do:true"`
	Id           any         //
	Name         any         //
	HostId       any         //
	Type         any         //
	LocalAddress any         //
	LocalPort    any         //
	RemoteHost   any         //
	RemotePort   any         //
	SocksHost    any         //
	SocksPort    any         //
	Enabled      any         //
	SyncUser     any         //
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
