// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SyncHostGroupsDao is the data access object for the table sync_host_groups.
type SyncHostGroupsDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  SyncHostGroupsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// SyncHostGroupsColumns defines and stores column names for the table sync_host_groups.
type SyncHostGroupsColumns struct {
	Id        string //
	Name      string //
	ParentId  string //
	SortOrder string //
	SyncUser  string //
	CreatedAt string //
	UpdatedAt string //
}

// syncHostGroupsColumns holds the columns for the table sync_host_groups.
var syncHostGroupsColumns = SyncHostGroupsColumns{
	Id:        "id",
	Name:      "name",
	ParentId:  "parent_id",
	SortOrder: "sort_order",
	SyncUser:  "sync_user",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewSyncHostGroupsDao creates and returns a new DAO object for table data access.
func NewSyncHostGroupsDao(handlers ...gdb.ModelHandler) *SyncHostGroupsDao {
	return &SyncHostGroupsDao{
		group:    "default",
		table:    "sync_host_groups",
		columns:  syncHostGroupsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SyncHostGroupsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SyncHostGroupsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SyncHostGroupsDao) Columns() SyncHostGroupsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SyncHostGroupsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SyncHostGroupsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SyncHostGroupsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
