package dao

import (
	"context"

	"ydsterm/internal/dbcore"
	"ydsterm/internal/dbcore/dao/internal"
	"ydsterm/internal/dbcore/model/entity"
	"ydsterm/internal/types"
)

type ydstermKeysDao struct {
	*internal.YdstermKeysDao
}

var (
	YdstermKeys = ydstermKeysDao{internal.NewYdstermKeysDao()}
)

func (d *ydstermKeysDao) Create(ctx context.Context, input types.KeyCreateInput) (*entity.YdstermKeys, error) {
	id := types.NewID()
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).Data(map[string]interface{}{
		d.Columns().Id:            id,
		d.Columns().Name:          input.Name,
		d.Columns().PrivateKeyEnc: "",
		d.Columns().PublicKey:     input.PublicKey,
		d.Columns().PassphraseEnc: "",
		d.Columns().SyncUser:      su,
	}).Insert()
	if err != nil {
		return nil, err
	}
	var key *entity.YdstermKeys
	err = d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Scan(&key)
	return key, err
}

func (d *ydstermKeysDao) Update(ctx context.Context, id string, data map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Update(data)
	return err
}

func (d *ydstermKeysDao) Delete(ctx context.Context, id string) error {
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Delete()
	return err
}

func (d *ydstermKeysDao) Get(ctx context.Context, id string) (*entity.YdstermKeys, error) {
	var key *entity.YdstermKeys
	su := dbcore.GetSyncUser(ctx)
	err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Scan(&key)
	return key, err
}

func (d *ydstermKeysDao) List(ctx context.Context) ([]entity.YdstermKeys, error) {
	su := dbcore.GetSyncUser(ctx)
	var keys []entity.YdstermKeys
	err := d.Ctx(ctx).Where(d.Columns().SyncUser, su).OrderAsc(d.Columns().Name).Scan(&keys)
	if err != nil {
		return nil, err
	}
	if keys == nil {
		keys = []entity.YdstermKeys{}
	}
	return keys, nil
}
