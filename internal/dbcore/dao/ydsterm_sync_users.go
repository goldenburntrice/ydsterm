package dao

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"ydsterm/internal/dbcore/dao/internal"
	"ydsterm/internal/types"
)

type ydstermSyncUsersDao struct {
	*internal.YdstermSyncUsersDao
}

var YdstermSyncUsers = ydstermSyncUsersDao{internal.NewYdstermSyncUsersDao()}

type SyncUserEntity struct {
	Id           string      `json:"id"           orm:"id"            description:""`
	Username     string      `json:"username"     orm:"username"      description:""`
	PasswordHash string      `json:"passwordHash" orm:"password_hash" description:""`
	IsActive     int         `json:"isActive"     orm:"is_active"     description:""`
	LastSyncAt   *gtime.Time `json:"lastSyncAt"   orm:"last_sync_at"  description:""`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:""`
}

func (e *SyncUserEntity) String() string {
	return fmt.Sprintf("SyncUser{id:%s, username:%s, active:%d}", e.Id, e.Username, e.IsActive)
}

func (d *ydstermSyncUsersDao) GetActive(ctx context.Context) (*SyncUserEntity, error) {
	var u *SyncUserEntity
	err := d.Ctx(ctx).Where(d.Columns().IsActive, 1).Scan(&u)
	return u, err
}

func (d *ydstermSyncUsersDao) SetActive(ctx context.Context, id string) error {
	return d.Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Exec("UPDATE "+d.Table()+" SET "+d.Columns().IsActive+" = 0")
		if err != nil {
			return err
		}
		_, err = tx.Model(d.Table()).Ctx(ctx).WherePri(id).Update(g.Map{d.Columns().IsActive: 1})
		return err
	})
}

func (d *ydstermSyncUsersDao) List(ctx context.Context) ([]SyncUserEntity, error) {
	var users []SyncUserEntity
	err := d.Ctx(ctx).OrderAsc(d.Columns().Username).Scan(&users)
	if err != nil {
		return nil, err
	}
	if users == nil {
		users = []SyncUserEntity{}
	}
	return users, nil
}

func (d *ydstermSyncUsersDao) Create(ctx context.Context, username, passwordHash string) (*SyncUserEntity, error) {
	id := types.NewID()
	now := gtime.Now()
	_, err := d.Ctx(ctx).Data(g.Map{
		d.Columns().Id:            id,
		d.Columns().Username:      username,
		d.Columns().PasswordHash:  passwordHash,
		d.Columns().IsActive:      0,
		d.Columns().CreatedAt:     now,
		d.Columns().UpdatedAt:     now,
	}).Insert()
	if err != nil {
		return nil, err
	}
	var u *SyncUserEntity
	err = d.Ctx(ctx).WherePri(id).Scan(&u)
	return u, err
}

func (d *ydstermSyncUsersDao) GetByUsername(ctx context.Context, username string) (*SyncUserEntity, error) {
	var u *SyncUserEntity
	err := d.Ctx(ctx).Where(d.Columns().Username, username).Scan(&u)
	return u, err
}

func (d *ydstermSyncUsersDao) UpdateLastSyncAt(ctx context.Context, id string, t *gtime.Time) error {
	_, err := d.Ctx(ctx).WherePri(id).Update(g.Map{
		d.Columns().LastSyncAt: t,
		d.Columns().UpdatedAt:  gtime.Now(),
	})
	return err
}

func (d *ydstermSyncUsersDao) UpdatePasswordHash(ctx context.Context, id, passwordHash string) error {
	_, err := d.Ctx(ctx).WherePri(id).Update(g.Map{
		d.Columns().PasswordHash: passwordHash,
		d.Columns().UpdatedAt:    gtime.Now(),
	})
	return err
}
