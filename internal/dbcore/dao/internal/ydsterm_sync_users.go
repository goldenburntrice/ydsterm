// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// YdstermSyncUsersDao is the data access object for the table ydsterm_sync_users.
type YdstermSyncUsersDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  YdstermSyncUsersColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// YdstermSyncUsersColumns defines and stores column names for the table ydsterm_sync_users.
type YdstermSyncUsersColumns struct {
	Id         string //
	Username   string //
	IsActive   string //
	LastSyncAt string //
	CreatedAt  string //
	UpdatedAt  string //
}

// ydstermSyncUsersColumns holds the columns for the table ydsterm_sync_users.
var ydstermSyncUsersColumns = YdstermSyncUsersColumns{
	Id:         "id",
	Username:   "username",
	IsActive:   "is_active",
	LastSyncAt: "last_sync_at",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewYdstermSyncUsersDao creates and returns a new DAO object for table data access.
func NewYdstermSyncUsersDao(handlers ...gdb.ModelHandler) *YdstermSyncUsersDao {
	return &YdstermSyncUsersDao{
		group:    "default",
		table:    "ydsterm_sync_users",
		columns:  ydstermSyncUsersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *YdstermSyncUsersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *YdstermSyncUsersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *YdstermSyncUsersDao) Columns() YdstermSyncUsersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *YdstermSyncUsersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *YdstermSyncUsersDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *YdstermSyncUsersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
