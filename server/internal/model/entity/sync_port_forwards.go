// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SyncPortForwards is the golang structure for table sync_port_forwards.
type SyncPortForwards struct {
	Id           string      `json:"id"           orm:"id"            ` //
	Name         string      `json:"name"         orm:"name"          ` //
	HostId       string      `json:"hostId"       orm:"host_id"       ` //
	Type         string      `json:"type"         orm:"type"          ` //
	LocalAddress string      `json:"localAddress" orm:"local_address" ` //
	LocalPort    int         `json:"localPort"    orm:"local_port"    ` //
	RemoteHost   string      `json:"remoteHost"   orm:"remote_host"   ` //
	RemotePort   int         `json:"remotePort"   orm:"remote_port"   ` //
	SocksHost    string      `json:"socksHost"    orm:"socks_host"    ` //
	SocksPort    int         `json:"socksPort"    orm:"socks_port"    ` //
	Enabled      int         `json:"enabled"      orm:"enabled"       ` //
	SyncUser     string      `json:"syncUser"     orm:"sync_user"     ` //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` //
}
