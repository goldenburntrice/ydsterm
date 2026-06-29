// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SyncSnippetsDao is the data access object for the table sync_snippets.
type SyncSnippetsDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  SyncSnippetsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// SyncSnippetsColumns defines and stores column names for the table sync_snippets.
type SyncSnippetsColumns struct {
	Id        string //
	Name      string //
	Content   string //
	Language  string //
	FolderId  string //
	SortOrder string //
	SyncUser  string //
	CreatedAt string //
	UpdatedAt string //
}

// syncSnippetsColumns holds the columns for the table sync_snippets.
var syncSnippetsColumns = SyncSnippetsColumns{
	Id:        "id",
	Name:      "name",
	Content:   "content",
	Language:  "language",
	FolderId:  "folder_id",
	SortOrder: "sort_order",
	SyncUser:  "sync_user",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewSyncSnippetsDao creates and returns a new DAO object for table data access.
func NewSyncSnippetsDao(handlers ...gdb.ModelHandler) *SyncSnippetsDao {
	return &SyncSnippetsDao{
		group:    "default",
		table:    "sync_snippets",
		columns:  syncSnippetsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SyncSnippetsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SyncSnippetsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SyncSnippetsDao) Columns() SyncSnippetsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SyncSnippetsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SyncSnippetsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SyncSnippetsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
