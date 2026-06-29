package dao

import (
	"context"

	"ydsterm/internal/dbcore"
	"ydsterm/internal/dbcore/dao/internal"
	"ydsterm/internal/dbcore/model/entity"
	"ydsterm/internal/types"
)

type ydstermPortForwardsDao struct {
	*internal.YdstermPortForwardsDao
}

var (
	YdstermPortForwards = ydstermPortForwardsDao{internal.NewYdstermPortForwardsDao()}
)

func (d *ydstermPortForwardsDao) Create(ctx context.Context, input types.PortForwardCreateInput) (*entity.YdstermPortForwards, error) {
	id := types.NewID()
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).Data(map[string]interface{}{
		d.Columns().Id:           id,
		d.Columns().Name:         input.Name,
		d.Columns().HostId:       input.HostID,
		d.Columns().Type:         input.Type,
		d.Columns().LocalAddress: input.LocalAddress,
		d.Columns().LocalPort:    input.LocalPort,
		d.Columns().RemoteHost:   input.RemoteHost,
		d.Columns().RemotePort:   input.RemotePort,
		d.Columns().SocksHost:    input.SocksHost,
		d.Columns().SocksPort:    input.SocksPort,
		d.Columns().Enabled:      0,
		d.Columns().SyncUser:     su,
	}).Insert()
	if err != nil {
		return nil, err
	}
	var pf *entity.YdstermPortForwards
	err = d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Scan(&pf)
	return pf, err
}

func (d *ydstermPortForwardsDao) Update(ctx context.Context, id string, data map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Update(data)
	return err
}

func (d *ydstermPortForwardsDao) Delete(ctx context.Context, id string) error {
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Delete()
	return err
}

func (d *ydstermPortForwardsDao) Get(ctx context.Context, id string) (*entity.YdstermPortForwards, error) {
	var pf *entity.YdstermPortForwards
	su := dbcore.GetSyncUser(ctx)
	err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Scan(&pf)
	return pf, err
}

func (d *ydstermPortForwardsDao) List(ctx context.Context, hostID string) ([]entity.YdstermPortForwards, error) {
	su := dbcore.GetSyncUser(ctx)
	m := d.Ctx(ctx).Where(d.Columns().SyncUser, su)
	if hostID != "" {
		m = m.Where(d.Columns().HostId, hostID)
	}
	var pfs []entity.YdstermPortForwards
	err := m.OrderAsc(d.Columns().Name).Scan(&pfs)
	if err != nil {
		return nil, err
	}
	if pfs == nil {
		pfs = []entity.YdstermPortForwards{}
	}
	return pfs, nil
}

func (d *ydstermPortForwardsDao) Toggle(ctx context.Context, id string, enabled bool) error {
	su := dbcore.GetSyncUser(ctx)
	val := 0
	if enabled {
		val = 1
	}
	_, err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Update(map[string]interface{}{
		d.Columns().Enabled: val,
	})
	return err
}
