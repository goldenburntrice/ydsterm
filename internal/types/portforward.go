package types

import "time"

type PortForward struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	HostID       string    `json:"hostId"`
	Type         string    `json:"type"`
	LocalAddress string    `json:"localAddress"`
	LocalPort    int       `json:"localPort"`
	RemoteHost   string    `json:"remoteHost"`
	RemotePort   int       `json:"remotePort"`
	SocksHost    string    `json:"socksHost"`
	SocksPort    int       `json:"socksPort"`
	Enabled      bool      `json:"enabled"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type PortForwardCreateInput struct {
	Name         string `json:"name"`
	HostID       string `json:"hostId"`
	Type         string `json:"type"`
	LocalAddress string `json:"localAddress,omitempty"`
	LocalPort    int    `json:"localPort,omitempty"`
	RemoteHost   string `json:"remoteHost,omitempty"`
	RemotePort   int    `json:"remotePort,omitempty"`
	SocksHost    string `json:"socksHost,omitempty"`
	SocksPort    int    `json:"socksPort,omitempty"`
}

type PortForwardUpdateInput struct {
	ID           string  `json:"id"`
	Name         *string `json:"name,omitempty"`
	HostID       *string `json:"hostId,omitempty"`
	Type         *string `json:"type,omitempty"`
	LocalAddress *string `json:"localAddress,omitempty"`
	LocalPort    *int    `json:"localPort,omitempty"`
	RemoteHost   *string `json:"remoteHost,omitempty"`
	RemotePort   *int    `json:"remotePort,omitempty"`
	SocksHost    *string `json:"socksHost,omitempty"`
	SocksPort    *int    `json:"socksPort,omitempty"`
}
