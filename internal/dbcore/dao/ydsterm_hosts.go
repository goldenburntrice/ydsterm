package dao

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"

	"ydsterm/internal/dbcore"
	"ydsterm/internal/dbcore/dao/internal"
	"ydsterm/internal/dbcore/model/entity"
	"ydsterm/internal/types"
)

type ydstermHostsDao struct {
	*internal.YdstermHostsDao
}

var (
	YdstermHosts = ydstermHostsDao{internal.NewYdstermHostsDao()}
)

func (d *ydstermHostsDao) Create(ctx context.Context, input types.HostCreateInput) (*entity.YdstermHosts, error) {
	id := types.NewID()
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).Data(map[string]interface{}{
		d.Columns().Id:          id,
		d.Columns().Name:        input.Name,
		d.Columns().Hostname:    input.Hostname,
		d.Columns().Port:        input.Port,
		d.Columns().Username:    input.Username,
		d.Columns().AuthMethod:  input.AuthMethod,
		d.Columns().PasswordEnc: "",
		d.Columns().KeyId:       input.KeyID,
		d.Columns().GroupId:     input.GroupID,
		d.Columns().Color:       input.Color,
		d.Columns().SyncUser:    su,
	}).Insert()
	if err != nil {
		return nil, err
	}
	var host *entity.YdstermHosts
	err = d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Scan(&host)
	return host, err
}

func (d *ydstermHostsDao) Update(ctx context.Context, id string, data map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Update(data)
	return err
}

func (d *ydstermHostsDao) Delete(ctx context.Context, id string) error {
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Update(map[string]interface{}{
		d.Columns().DeletedAt: gtime.Now(),
	})
	return err
}

func (d *ydstermHostsDao) Get(ctx context.Context, id string) (*entity.YdstermHosts, error) {
	var host *entity.YdstermHosts
	su := dbcore.GetSyncUser(ctx)
	err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Scan(&host)
	return host, err
}

func (d *ydstermHostsDao) List(ctx context.Context, groupID string) ([]entity.YdstermHosts, error) {
	su := dbcore.GetSyncUser(ctx)
	m := d.Ctx(ctx).WhereNull(d.Columns().DeletedAt).Where(d.Columns().SyncUser, su)
	if groupID != "" {
		m = m.Where(d.Columns().GroupId, groupID)
	}
	var hosts []entity.YdstermHosts
	err := m.OrderAsc(d.Columns().Name).Scan(&hosts)
	if err != nil {
		return nil, err
	}
	if hosts == nil {
		hosts = []entity.YdstermHosts{}
	}
	return hosts, nil
}

func (d *ydstermHostsDao) ListAllForSync(ctx context.Context) ([]entity.YdstermHosts, error) {
	var hosts []entity.YdstermHosts
	err := d.Ctx(ctx).WhereNull(d.Columns().DeletedAt).OrderAsc(d.Columns().Name).Scan(&hosts)
	if err != nil {
		return nil, err
	}
	if hosts == nil {
		hosts = []entity.YdstermHosts{}
	}
	return hosts, nil
}

func (d *ydstermHostsDao) ListChangedSince(ctx context.Context, since string) ([]entity.YdstermHosts, error) {
	su := dbcore.GetSyncUser(ctx)
	var hosts []entity.YdstermHosts
	m := d.Ctx(ctx).WhereNull(d.Columns().DeletedAt).Where(d.Columns().SyncUser, su)
	if since != "" {
		m = m.WhereGT(d.Columns().UpdatedAt, since)
	}
	err := m.OrderAsc(d.Columns().Name).Scan(&hosts)
	if err != nil {
		return nil, err
	}
	if hosts == nil {
		hosts = []entity.YdstermHosts{}
	}
	return hosts, nil
}
