// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SyncPortForwardsDao is the data access object for the table sync_port_forwards.
type SyncPortForwardsDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  SyncPortForwardsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// SyncPortForwardsColumns defines and stores column names for the table sync_port_forwards.
type SyncPortForwardsColumns struct {
	Id           string //
	Name         string //
	HostId       string //
	Type         string //
	LocalAddress string //
	LocalPort    string //
	RemoteHost   string //
	RemotePort   string //
	SocksHost    string //
	SocksPort    string //
	Enabled      string //
	SyncUser     string //
	CreatedAt    string //
	UpdatedAt    string //
}

// syncPortForwardsColumns holds the columns for the table sync_port_forwards.
var syncPortForwardsColumns = SyncPortForwardsColumns{
	Id:           "id",
	Name:         "name",
	HostId:       "host_id",
	Type:         "type",
	LocalAddress: "local_address",
	LocalPort:    "local_port",
	RemoteHost:   "remote_host",
	RemotePort:   "remote_port",
	SocksHost:    "socks_host",
	SocksPort:    "socks_port",
	Enabled:      "enabled",
	SyncUser:     "sync_user",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewSyncPortForwardsDao creates and returns a new DAO object for table data access.
func NewSyncPortForwardsDao(handlers ...gdb.ModelHandler) *SyncPortForwardsDao {
	return &SyncPortForwardsDao{
		group:    "default",
		table:    "sync_port_forwards",
		columns:  syncPortForwardsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SyncPortForwardsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SyncPortForwardsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SyncPortForwardsDao) Columns() SyncPortForwardsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SyncPortForwardsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SyncPortForwardsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *SyncPortForwardsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
