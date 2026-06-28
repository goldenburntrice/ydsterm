// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// YdstermKeysDao is the data access object for the table ydsterm_keys.
type YdstermKeysDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  YdstermKeysColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// YdstermKeysColumns defines and stores column names for the table ydsterm_keys.
type YdstermKeysColumns struct {
	Id            string //
	Name          string //
	PrivateKeyEnc string //
	PublicKey     string //
	PassphraseEnc string //
	CreatedAt     string //
	UpdatedAt     string //
}

// ydstermKeysColumns holds the columns for the table ydsterm_keys.
var ydstermKeysColumns = YdstermKeysColumns{
	Id:            "id",
	Name:          "name",
	PrivateKeyEnc: "private_key_enc",
	PublicKey:     "public_key",
	PassphraseEnc: "passphrase_enc",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewYdstermKeysDao creates and returns a new DAO object for table data access.
func NewYdstermKeysDao(handlers ...gdb.ModelHandler) *YdstermKeysDao {
	return &YdstermKeysDao{
		group:    "default",
		table:    "ydsterm_keys",
		columns:  ydstermKeysColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *YdstermKeysDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *YdstermKeysDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *YdstermKeysDao) Columns() YdstermKeysColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *YdstermKeysDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *YdstermKeysDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *YdstermKeysDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
