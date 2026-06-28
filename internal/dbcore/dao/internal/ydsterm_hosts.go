// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// YdstermHostsDao is the data access object for the table ydsterm_hosts.
type YdstermHostsDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  YdstermHostsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// YdstermHostsColumns defines and stores column names for the table ydsterm_hosts.
type YdstermHostsColumns struct {
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
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
}

// ydstermHostsColumns holds the columns for the table ydsterm_hosts.
var ydstermHostsColumns = YdstermHostsColumns{
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
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewYdstermHostsDao creates and returns a new DAO object for table data access.
func NewYdstermHostsDao(handlers ...gdb.ModelHandler) *YdstermHostsDao {
	return &YdstermHostsDao{
		group:    "default",
		table:    "ydsterm_hosts",
		columns:  ydstermHostsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *YdstermHostsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *YdstermHostsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *YdstermHostsDao) Columns() YdstermHostsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *YdstermHostsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *YdstermHostsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *YdstermHostsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
