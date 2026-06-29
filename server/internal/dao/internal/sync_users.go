// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SyncUsersDao is the data access object for the table sync_users.
type SyncUsersDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SyncUsersColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SyncUsersColumns defines and stores column names for the table sync_users.
type SyncUsersColumns struct {
	Id           string //
	Username     string //
	PasswordHash string //
	CreatedAt    string //
	UpdatedAt    string //
}

// syncUsersColumns holds the columns for the table sync_users.
var syncUsersColumns = SyncUsersColumns{
	Id:           "id",
	Username:     "username",
	PasswordHash: "password_hash",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewSyncUsersDao creates and returns a new DAO object for table data access.
func NewSyncUsersDao(handlers ...gdb.ModelHandler) *SyncUsersDao {
	return &SyncUsersDao{
		group:    "default",
		table:    "sync_users",
		columns:  syncUsersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SyncUsersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SyncUsersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SyncUsersDao) Columns() SyncUsersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SyncUsersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SyncUsersDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SyncUsersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
