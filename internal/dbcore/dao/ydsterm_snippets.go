package dao

import (
	"context"

	"ydsterm/internal/dbcore"
	"ydsterm/internal/dbcore/dao/internal"
	"ydsterm/internal/dbcore/model/entity"
	"ydsterm/internal/types"
)

type ydstermSnippetsDao struct {
	*internal.YdstermSnippetsDao
}

var (
	YdstermSnippets = ydstermSnippetsDao{internal.NewYdstermSnippetsDao()}
)

func (d *ydstermSnippetsDao) Create(ctx context.Context, input types.SnippetCreateInput) (*entity.YdstermSnippets, error) {
	id := types.NewID()
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).Data(map[string]interface{}{
		d.Columns().Id:       id,
		d.Columns().Name:     input.Name,
		d.Columns().Content:  input.Content,
		d.Columns().Language: input.Language,
		d.Columns().FolderId: input.FolderID,
		d.Columns().SyncUser: su,
	}).Insert()
	if err != nil {
		return nil, err
	}
	var s *entity.YdstermSnippets
	err = d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Scan(&s)
	return s, err
}

func (d *ydstermSnippetsDao) Update(ctx context.Context, id string, data map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Update(data)
	return err
}

func (d *ydstermSnippetsDao) Delete(ctx context.Context, id string) error {
	su := dbcore.GetSyncUser(ctx)
	_, err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Delete()
	return err
}

func (d *ydstermSnippetsDao) Get(ctx context.Context, id string) (*entity.YdstermSnippets, error) {
	var s *entity.YdstermSnippets
	su := dbcore.GetSyncUser(ctx)
	err := d.Ctx(ctx).WherePri(id).Where(d.Columns().SyncUser, su).Scan(&s)
	return s, err
}

func (d *ydstermSnippetsDao) List(ctx context.Context, folderID string) ([]entity.YdstermSnippets, error) {
	su := dbcore.GetSyncUser(ctx)
	m := d.Ctx(ctx).Where(d.Columns().SyncUser, su)
	if folderID != "" {
		m = m.Where(d.Columns().FolderId, folderID)
	}
	var snippets []entity.YdstermSnippets
	err := m.OrderAsc(d.Columns().SortOrder).OrderAsc(d.Columns().Name).Scan(&snippets)
	if err != nil {
		return nil, err
	}
	if snippets == nil {
		snippets = []entity.YdstermSnippets{}
	}
	return snippets, nil
}
