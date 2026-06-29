// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SyncHostsDao is the data access object for the table sync_hosts.
type SyncHostsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SyncHostsColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SyncHostsColumns defines and stores column names for the table sync_hosts.
type SyncHostsColumns struct {
	Id          string //
	Name        string //
	Hostname    string //
	Port        string //
	Username    string //
	AuthMethod  string //
	PasswordEnc string //
	KeyId       string //
	GroupId     string //
	Color       string //
	SyncUser    string //
	CreatedAt   string //
	UpdatedAt   string //
}

// syncHostsColumns holds the columns for the table sync_hosts.
var syncHostsColumns = SyncHostsColumns{
	Id:          "id",
	Name:        "name",
	Hostname:    "hostname",
	Port:        "port",
	Username:    "username",
	AuthMethod:  "auth_method",
	PasswordEnc: "password_enc",
	KeyId:       "key_id",
	GroupId:     "group_id",
	Color:       "color",
	SyncUser:    "sync_user",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewSyncHostsDao creates and returns a new DAO object for table data access.
func NewSyncHostsDao(handlers ...gdb.ModelHandler) *SyncHostsDao {
	return &SyncHostsDao{
		group:    "default",
		table:    "sync_hosts",
		columns:  syncHostsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SyncHostsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SyncHostsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SyncHostsDao) Columns() SyncHostsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SyncHostsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SyncHostsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SyncHostsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
