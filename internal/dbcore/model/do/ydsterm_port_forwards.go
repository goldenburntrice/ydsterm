// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermPortForwards is the golang structure of table ydsterm_port_forwards for DAO operations like Where/Data.
type YdstermPortForwards struct {
	g.Meta       `orm:"table:ydsterm_port_forwards, do:true"`
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
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
