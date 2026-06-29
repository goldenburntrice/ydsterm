package dao

import (
	"context"

	"ydsterm/internal/dbcore"
	"ydsterm/internal/dbcore/dao/internal"
	"ydsterm/internal/dbcore/model/entity"
	"ydsterm/internal/types"
)

type ydstermHostGroupsDao struct {
	*internal.YdstermHostGroupsDao
}

var (
	YdstermHostGroups = ydstermHostGroupsDao{internal.NewYdstermHostGroupsDao()}
)

func (d *ydstermHostGroupsDao) Create(ctx context.Context, input types.HostGroupCreateInput) (*entity.YdstermHostGroups, error) {
	id := types.NewID()
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).Data(map[string]interface{}{
		d.Columns().Id:       id,
		d.Columns().Name:     input.Name,
		d.Columns().ParentId: input.ParentID,
		d.Columns().SyncUser: su,
	}).Insert()
	if err != nil {
		return nil, err
	}
	var g *entity.YdstermHostGroups
	err = d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Scan(&g)
	return g, err
}

func (d *ydstermHostGroupsDao) Update(ctx context.Context, id string, data map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Update(data)
	return err
}

func (d *ydstermHostGroupsDao) Delete(ctx context.Context, id string) error {
	su := dbcore.GetSyncUser(ctx)
	_, _ = YdstermHosts.Ctx(ctx).Data(map[string]interface{}{
		YdstermHosts.Columns().GroupId: "",
	}).Where(YdstermHosts.Columns().GroupId, id).Where(YdstermHosts.Columns().SyncUser, su).Update()
	_, err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Delete()
	return err
}

func (d *ydstermHostGroupsDao) Get(ctx context.Context, id string) (*entity.YdstermHostGroups, error) {
	var g *entity.YdstermHostGroups
	su := dbcore.GetSyncUser(ctx)
	err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Scan(&g)
	return g, err
}

func (d *ydstermHostGroupsDao) List(ctx context.Context) ([]entity.YdstermHostGroups, error) {
	su := dbcore.GetSyncUser(ctx)
	var groups []entity.YdstermHostGroups
	err := d.Ctx(ctx).Where(d.Columns().SyncUser, su).OrderAsc(d.Columns().SortOrder).OrderAsc(d.Columns().Name).Scan(&groups)
	if err != nil {
		return nil, err
	}
	if groups == nil {
		groups = []entity.YdstermHostGroups{}
	}
	return groups, nil
}
