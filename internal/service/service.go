package service

import "ydsterm/internal/types"

type IHostService interface {
	Create(input types.HostCreateInput) (*types.Host, error)
	Update(input types.HostUpdateInput) (*types.Host, error)
	Delete(id string) error
	Get(id string) (*types.Host, error)
	List(groupID string) ([]types.Host, error)
	CreateKey(input types.KeyCreateInput) (*types.Key, error)
	UpdateKey(input types.KeyUpdateInput) (*types.Key, error)
	DeleteKey(id string) error
	GetKey(id string) (*types.Key, error)
	ListKeys() ([]types.Key, error)
	CreateGroup(input types.HostGroupCreateInput) (*types.HostGroup, error)
	UpdateGroup(input types.HostGroupUpdateInput) (*types.HostGroup, error)
	DeleteGroup(id string) error
	ListGroups() ([]types.HostGroup, error)
}

type ISnippetService interface {
	Create(input types.SnippetCreateInput) (*types.Snippet, error)
	Update(input types.SnippetUpdateInput) (*types.Snippet, error)
	Delete(id string) error
	Get(id string) (*types.Snippet, error)
	List(folderID string) ([]types.Snippet, error)
	CreateFolder(input types.SnippetFolderCreateInput) (*types.SnippetFolder, error)
	UpdateFolder(input types.SnippetFolderUpdateInput) (*types.SnippetFolder, error)
	DeleteFolder(id string) error
	ListFolders() ([]types.SnippetFolder, error)
}

type IPortForwardService interface {
	Create(input types.PortForwardCreateInput) (*types.PortForward, error)
	Update(input types.PortForwardUpdateInput) (*types.PortForward, error)
	Delete(id string) error
	Get(id string) (*types.PortForward, error)
	List(hostID string) ([]types.PortForward, error)
	Toggle(id string, enabled bool) error
}

type ISettingsService interface {
	Get(key string) (string, error)
	Set(key, value string) error
	VerifyUser(serverAddr, serverKey, username, password string) (*types.VerifyResponse, error)
}

type IUserService interface {
	GetCurrentUser() (*types.SyncUser, error)
	ListUsers() ([]types.SyncUser, error)
	SwitchUser(userID string) error
	RegisterVerifiedUser(username string, serverUserID string) (*types.SyncUser, error)
	InitActiveUser() error
}
