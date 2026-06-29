package dao

import (
	"context"

	"ydsterm/internal/dbcore"
	"ydsterm/internal/dbcore/dao/internal"
	"ydsterm/internal/dbcore/model/entity"
	"ydsterm/internal/types"
)

type ydstermSnippetFoldersDao struct {
	*internal.YdstermSnippetFoldersDao
}

var (
	YdstermSnippetFolders = ydstermSnippetFoldersDao{internal.NewYdstermSnippetFoldersDao()}
)

func (d *ydstermSnippetFoldersDao) Create(ctx context.Context, input types.SnippetFolderCreateInput) (*entity.YdstermSnippetFolders, error) {
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
	var f *entity.YdstermSnippetFolders
	err = d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Scan(&f)
	return f, err
}

func (d *ydstermSnippetFoldersDao) Update(ctx context.Context, id string, data map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Update(data)
	return err
}

func (d *ydstermSnippetFoldersDao) Delete(ctx context.Context, id string) error {
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Delete()
	return err
}

func (d *ydstermSnippetFoldersDao) Get(ctx context.Context, id string) (*entity.YdstermSnippetFolders, error) {
	var f *entity.YdstermSnippetFolders
	su := dbcore.GetSyncUser(ctx)
	err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Scan(&f)
	return f, err
}

func (d *ydstermSnippetFoldersDao) List(ctx context.Context) ([]entity.YdstermSnippetFolders, error) {
	su := dbcore.GetSyncUser(ctx)
	var folders []entity.YdstermSnippetFolders
	err := d.Ctx(ctx).Where(d.Columns().SyncUser, su).OrderAsc(d.Columns().SortOrder).OrderAsc(d.Columns().Name).Scan(&folders)
	if err != nil {
		return nil, err
	}
	if folders == nil {
		folders = []entity.YdstermSnippetFolders{}
	}
	return folders, nil
}
