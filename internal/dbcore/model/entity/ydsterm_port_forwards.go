// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermPortForwards is the golang structure for table ydsterm_port_forwards.
type YdstermPortForwards struct {
	Id           string      `json:"id"           orm:"id"            description:""`
	Name         string      `json:"name"         orm:"name"          description:""`
	HostId       string      `json:"hostId"       orm:"host_id"       description:""`
	Type         string      `json:"type"         orm:"type"          description:""`
	LocalAddress string      `json:"localAddress" orm:"local_address" description:""`
	LocalPort    int         `json:"localPort"    orm:"local_port"    description:""`
	RemoteHost   string      `json:"remoteHost"   orm:"remote_host"   description:""`
	RemotePort   int         `json:"remotePort"   orm:"remote_port"   description:""`
	SocksHost    string      `json:"socksHost"    orm:"socks_host"    description:""`
	SocksPort    int         `json:"socksPort"    orm:"socks_port"    description:""`
	Enabled      int         `json:"enabled"      orm:"enabled"       description:""`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:""`
}
