package service

import (
	"context"

	"ydsterm/internal/dbcore/dao"
	"ydsterm/internal/types"
)

type SnippetServiceImpl struct{}

func NewSnippetService() *SnippetServiceImpl { return &SnippetServiceImpl{} }

func (s *SnippetServiceImpl) Create(input types.SnippetCreateInput) (*types.Snippet, error) {
	snip, err := dao.YdstermSnippets.Create(context.Background(), input)
	if err != nil {
		return nil, err
	}
	return entityToSnippet(snip), nil
}

func (s *SnippetServiceImpl) Update(input types.SnippetUpdateInput) (*types.Snippet, error) {
	ctx := context.Background()
	data := make(map[string]interface{})
	cols := dao.YdstermSnippets.Columns()

	if input.Name != nil {
		data[cols.Name] = *input.Name
	}
	if input.Content != nil {
		data[cols.Content] = *input.Content
	}
	if input.Language != nil {
		data[cols.Language] = *input.Language
	}
	if input.FolderID != nil {
		data[cols.FolderId] = *input.FolderID
	}

	if err := dao.YdstermSnippets.Update(ctx, input.ID, data); err != nil {
		return nil, err
	}
	return s.Get(input.ID)
}

func (s *SnippetServiceImpl) Delete(id string) error {
	return dao.YdstermSnippets.Delete(context.Background(), id)
}

func (s *SnippetServiceImpl) Get(id string) (*types.Snippet, error) {
	snip, err := dao.YdstermSnippets.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return entityToSnippet(snip), nil
}

func (s *SnippetServiceImpl) List(folderID string) ([]types.Snippet, error) {
	snippets, err := dao.YdstermSnippets.List(context.Background(), folderID)
	if err != nil {
		return nil, err
	}
	out := make([]types.Snippet, len(snippets))
	for i, s := range snippets {
		out[i] = *entityToSnippet(&s)
	}
	return out, nil
}

func (s *SnippetServiceImpl) CreateFolder(input types.SnippetFolderCreateInput) (*types.SnippetFolder, error) {
	f, err := dao.YdstermSnippetFolders.Create(context.Background(), input)
	if err != nil {
		return nil, err
	}
	return entityToSnippetFolder(f), nil
}

func (s *SnippetServiceImpl) UpdateFolder(input types.SnippetFolderUpdateInput) (*types.SnippetFolder, error) {
	ctx := context.Background()
	data := make(map[string]interface{})
	cols := dao.YdstermSnippetFolders.Columns()
	if input.Name != nil {
		data[cols.Name] = *input.Name
	}
	if input.ParentID != nil {
		data[cols.ParentId] = *input.ParentID
	}
	if err := dao.YdstermSnippetFolders.Update(ctx, input.ID, data); err != nil {
		return nil, err
	}
	f, err := dao.YdstermSnippetFolders.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return entityToSnippetFolder(f), nil
}

func (s *SnippetServiceImpl) DeleteFolder(id string) error {
	return dao.YdstermSnippetFolders.Delete(context.Background(), id)
}

func (s *SnippetServiceImpl) ListFolders() ([]types.SnippetFolder, error) {
	folders, err := dao.YdstermSnippetFolders.List(context.Background())
	if err != nil {
		return nil, err
	}
	out := make([]types.SnippetFolder, len(folders))
	for i, f := range folders {
		out[i] = *entityToSnippetFolder(&f)
	}
	return out, nil
}
