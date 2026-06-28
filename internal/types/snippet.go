package types

import "time"

type Snippet struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Language  string    `json:"language"`
	FolderID  string    `json:"folderId"`
	SortOrder int       `json:"sortOrder"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SnippetCreateInput struct {
	Name     string `json:"name"`
	Content  string `json:"content"`
	Language string `json:"language"`
	FolderID string `json:"folderId"`
}

type SnippetUpdateInput struct {
	ID       string  `json:"id"`
	Name     *string `json:"name,omitempty"`
	Content  *string `json:"content,omitempty"`
	Language *string `json:"language,omitempty"`
	FolderID *string `json:"folderId,omitempty"`
}

type SnippetFolder struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ParentID  string    `json:"parentId"`
	SortOrder int       `json:"sortOrder"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SnippetFolderCreateInput struct {
	Name     string `json:"name"`
	ParentID string `json:"parentId"`
}

type SnippetFolderUpdateInput struct {
	ID       string  `json:"id"`
	Name     *string `json:"name,omitempty"`
	ParentID *string `json:"parentId,omitempty"`
}
