package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"ydsterm/internal/dbcore/dao/internal"
)

type ydstermSettingsDao struct {
	*internal.YdstermSettingsDao
}

var YdstermSettings = ydstermSettingsDao{internal.NewYdstermSettingsDao()}

func (d *ydstermSettingsDao) Get(ctx context.Context, key string) (string, error) {
	val, err := d.Ctx(ctx).Where(d.Columns().Key, key).Value(d.Columns().Value)
	if err != nil {
		return "", err
	}
	return val.String(), nil
}

func (d *ydstermSettingsDao) Set(ctx context.Context, key, value string) error {
	now := gtime.Now()
	existing, err := d.Ctx(ctx).Where(d.Columns().Key, key).One()
	if err != nil {
		return err
	}
	if existing.IsEmpty() {
		_, err = d.Ctx(ctx).Data(g.Map{
			d.Columns().Key:       key,
			d.Columns().Value:     value,
			d.Columns().UpdatedAt: now,
		}).Insert()
	} else {
		_, err = d.Ctx(ctx).Where(d.Columns().Key, key).Update(g.Map{
			d.Columns().Value:     value,
			d.Columns().UpdatedAt: now,
		})
	}
	return err
}

func (d *ydstermSettingsDao) Delete(ctx context.Context, key string) error {
	_, err := d.Ctx(ctx).Where(d.Columns().Key, key).Delete()
	return err
}
