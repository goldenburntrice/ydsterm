package service

import (
	"time"

	"github.com/gogf/gf/v2/os/gtime"

	"ydsterm/internal/dbcore/model/entity"
	"ydsterm/internal/types"
)

func entityToHost(e *entity.YdstermHosts) *types.Host {
	if e == nil {
		return nil
	}
	return &types.Host{
		ID:         e.Id,
		Name:       e.Name,
		Hostname:   e.Hostname,
		Port:       e.Port,
		Username:   e.Username,
		AuthMethod: e.AuthMethod,
		KeyID:      e.KeyId,
		GroupID:    e.GroupId,
		Color:      e.Color,
		CreatedAt:  gtimeToTime(e.CreatedAt),
		UpdatedAt:  gtimeToTime(e.UpdatedAt),
	}
}

func entityToKey(e *entity.YdstermKeys) *types.Key {
	if e == nil {
		return nil
	}
	return &types.Key{
		ID:        e.Id,
		Name:      e.Name,
		PublicKey: e.PublicKey,
		CreatedAt: gtimeToTime(e.CreatedAt),
		UpdatedAt: gtimeToTime(e.UpdatedAt),
	}
}

func entityToHostGroup(e *entity.YdstermHostGroups) *types.HostGroup {
	if e == nil {
		return nil
	}
	return &types.HostGroup{
		ID:        e.Id,
		Name:      e.Name,
		ParentID:  e.ParentId,
		SortOrder: e.SortOrder,
		CreatedAt: gtimeToTime(e.CreatedAt),
		UpdatedAt: gtimeToTime(e.UpdatedAt),
	}
}

func entityToSnippet(e *entity.YdstermSnippets) *types.Snippet {
	if e == nil {
		return nil
	}
	return &types.Snippet{
		ID:        e.Id,
		Name:      e.Name,
		Content:   e.Content,
		Language:  e.Language,
		FolderID:  e.FolderId,
		SortOrder: e.SortOrder,
		CreatedAt: gtimeToTime(e.CreatedAt),
		UpdatedAt: gtimeToTime(e.UpdatedAt),
	}
}

func entityToSnippetFolder(e *entity.YdstermSnippetFolders) *types.SnippetFolder {
	if e == nil {
		return nil
	}
	return &types.SnippetFolder{
		ID:        e.Id,
		Name:      e.Name,
		ParentID:  e.ParentId,
		SortOrder: e.SortOrder,
		CreatedAt: gtimeToTime(e.CreatedAt),
		UpdatedAt: gtimeToTime(e.UpdatedAt),
	}
}

func entityToPortForward(e *entity.YdstermPortForwards) *types.PortForward {
	if e == nil {
		return nil
	}
	return &types.PortForward{
		ID:           e.Id,
		Name:         e.Name,
		HostID:       e.HostId,
		Type:         e.Type,
		LocalAddress: e.LocalAddress,
		LocalPort:    e.LocalPort,
		RemoteHost:   e.RemoteHost,
		RemotePort:   e.RemotePort,
		SocksHost:    e.SocksHost,
		SocksPort:    e.SocksPort,
		Enabled:      e.Enabled == 1,
		CreatedAt:    gtimeToTime(e.CreatedAt),
		UpdatedAt:    gtimeToTime(e.UpdatedAt),
	}
}

func gtimeToTime(gt *gtime.Time) time.Time {
	if gt == nil {
		return time.Time{}
	}
	return gt.Time
}
