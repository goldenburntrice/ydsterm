package service

import (
	"context"

	"ydsterm/internal/dbcore"
	"ydsterm/internal/dbcore/dao"
	"ydsterm/internal/dbcore/model/entity"
	"ydsterm/internal/crypto"
	"ydsterm/internal/types"
)

func activeCtx() context.Context {
	return dbcore.WithActiveUser(context.Background())
}

type HostServiceImpl struct{}

func NewHostService() *HostServiceImpl { return &HostServiceImpl{} }

var syncService *SyncServiceImpl

func SetSyncService(s *SyncServiceImpl) {
	syncService = s
}

func (s *HostServiceImpl) Create(input types.HostCreateInput) (*types.Host, error) {
	ctx := activeCtx()
	passwordEnc := ""
	if input.Password != "" {
		var err error
		passwordEnc, err = crypto.Encrypt(input.Password)
		if err != nil {
			return nil, err
		}
	}
	input.Password = ""
	host, err := dao.YdstermHosts.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	if passwordEnc != "" {
		_ = dao.YdstermHosts.Update(ctx, host.Id, map[string]interface{}{
			dao.YdstermHosts.Columns().PasswordEnc: passwordEnc,
		})
		host.PasswordEnc = passwordEnc
	}
	if syncService != nil {
		go syncService.PushHost(entityToMap(host))
	}
	return entityToHost(host), nil
}

func (s *HostServiceImpl) Update(input types.HostUpdateInput) (*types.Host, error) {
	ctx := activeCtx()
	data := make(map[string]interface{})
	cols := dao.YdstermHosts.Columns()

	if input.Name != nil {
		data[cols.Name] = *input.Name
	}
	if input.Hostname != nil {
		data[cols.Hostname] = *input.Hostname
	}
	if input.Port != nil {
		data[cols.Port] = *input.Port
	}
	if input.Username != nil {
		data[cols.Username] = *input.Username
	}
	if input.AuthMethod != nil {
		data[cols.AuthMethod] = *input.AuthMethod
	}
	if input.Password != nil {
		enc, err := crypto.Encrypt(*input.Password)
		if err != nil {
			return nil, err
		}
		data[cols.PasswordEnc] = enc
	}
	if input.KeyID != nil {
		data[cols.KeyId] = *input.KeyID
	}
	if input.GroupID != nil {
		data[cols.GroupId] = *input.GroupID
	}
	if input.Color != nil {
		data[cols.Color] = *input.Color
	}

	if err := dao.YdstermHosts.Update(ctx, input.ID, data); err != nil {
		return nil, err
	}
	host, _ := dao.YdstermHosts.Get(ctx, input.ID)
	if syncService != nil && host != nil {
		go syncService.PushHost(entityToMap(host))
	}
	return entityToHost(host), nil
}

func (s *HostServiceImpl) Delete(id string) error {
	if err := dao.YdstermHosts.Delete(activeCtx(), id); err != nil {
		return err
	}
	if syncService != nil {
		go syncService.DeleteHost(id)
	}
	return nil
}

func (s *HostServiceImpl) Get(id string) (*types.Host, error) {
	host, err := dao.YdstermHosts.Get(activeCtx(), id)
	if err != nil {
		return nil, err
	}
	return entityToHost(host), nil
}

func (s *HostServiceImpl) List(groupID string) ([]types.Host, error) {
	hosts, err := dao.YdstermHosts.List(activeCtx(), groupID)
	if err != nil {
		return nil, err
	}
	out := make([]types.Host, len(hosts))
	for i, h := range hosts {
		out[i] = *entityToHost(&h)
	}
	return out, nil
}

func (s *HostServiceImpl) CreateKey(input types.KeyCreateInput) (*types.Key, error) {
	ctx := activeCtx()
	privEnc, err := crypto.Encrypt(input.PrivateKey)
	if err != nil {
		return nil, err
	}
	passEnc := ""
	if input.Passphrase != "" {
		passEnc, err = crypto.Encrypt(input.Passphrase)
		if err != nil {
			return nil, err
		}
	}
	input.PrivateKey = ""
	input.Passphrase = ""
	key, err := dao.YdstermKeys.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	cols := dao.YdstermKeys.Columns()
	upd := map[string]interface{}{cols.PrivateKeyEnc: privEnc}
	if passEnc != "" {
		upd[cols.PassphraseEnc] = passEnc
	}
	_ = dao.YdstermKeys.Update(ctx, key.Id, upd)
	key.PrivateKeyEnc = privEnc
	key.PassphraseEnc = passEnc
	if syncService != nil {
		go syncService.PushKey(keyEntityToMap(key))
	}
	return entityToKey(key), nil
}

func (s *HostServiceImpl) UpdateKey(input types.KeyUpdateInput) (*types.Key, error) {
	ctx := activeCtx()
	data := make(map[string]interface{})
	cols := dao.YdstermKeys.Columns()

	if input.Name != nil {
		data[cols.Name] = *input.Name
	}
	if input.PrivateKey != nil {
		enc, err := crypto.Encrypt(*input.PrivateKey)
		if err != nil {
			return nil, err
		}
		data[cols.PrivateKeyEnc] = enc
	}
	if input.PublicKey != nil {
		data[cols.PublicKey] = *input.PublicKey
	}
	if input.Passphrase != nil {
		enc, err := crypto.Encrypt(*input.Passphrase)
		if err != nil {
			return nil, err
		}
		data[cols.PassphraseEnc] = enc
	}

	if err := dao.YdstermKeys.Update(ctx, input.ID, data); err != nil {
		return nil, err
	}
	key, _ := dao.YdstermKeys.Get(ctx, input.ID)
	if syncService != nil && key != nil {
		go syncService.PushKey(keyEntityToMap(key))
	}
	return entityToKey(key), nil
}

func (s *HostServiceImpl) DeleteKey(id string) error {
	if err := dao.YdstermKeys.Delete(activeCtx(), id); err != nil {
		return err
	}
	if syncService != nil {
		go syncService.DeleteKey(id)
	}
	return nil
}

func (s *HostServiceImpl) GetKey(id string) (*types.Key, error) {
	key, err := dao.YdstermKeys.Get(activeCtx(), id)
	if err != nil {
		return nil, err
	}
	return entityToKey(key), nil
}

func (s *HostServiceImpl) ListKeys() ([]types.Key, error) {
	keys, err := dao.YdstermKeys.List(activeCtx())
	if err != nil {
		return nil, err
	}
	out := make([]types.Key, len(keys))
	for i, k := range keys {
		out[i] = *entityToKey(&k)
	}
	return out, nil
}

func (s *HostServiceImpl) CreateGroup(input types.HostGroupCreateInput) (*types.HostGroup, error) {
	g, err := dao.YdstermHostGroups.Create(activeCtx(), input)
	if err != nil {
		return nil, err
	}
	if syncService != nil {
		go syncService.PushHostGroup(groupEntityToMap(g))
	}
	return entityToHostGroup(g), nil
}

func (s *HostServiceImpl) UpdateGroup(input types.HostGroupUpdateInput) (*types.HostGroup, error) {
	ctx := activeCtx()
	data := make(map[string]interface{})
	cols := dao.YdstermHostGroups.Columns()
	if input.Name != nil {
		data[cols.Name] = *input.Name
	}
	if input.ParentID != nil {
		data[cols.ParentId] = *input.ParentID
	}
	if err := dao.YdstermHostGroups.Update(ctx, input.ID, data); err != nil {
		return nil, err
	}
	g, _ := dao.YdstermHostGroups.Get(ctx, input.ID)
	if syncService != nil && g != nil {
		go syncService.PushHostGroup(groupEntityToMap(g))
	}
	return entityToHostGroup(g), nil
}

func (s *HostServiceImpl) DeleteGroup(id string) error {
	if err := dao.YdstermHostGroups.Delete(activeCtx(), id); err != nil {
		return err
	}
	if syncService != nil {
		go syncService.DeleteHostGroup(id)
	}
	return nil
}

func (s *HostServiceImpl) ListGroups() ([]types.HostGroup, error) {
	groups, err := dao.YdstermHostGroups.List(activeCtx())
	if err != nil {
		return nil, err
	}
	out := make([]types.HostGroup, len(groups))
	for i, g := range groups {
		out[i] = *entityToHostGroup(&g)
	}
	return out, nil
}

func entityToMap(e *entity.YdstermHosts) map[string]interface{} {
	if e == nil {
		return nil
	}
	m := map[string]interface{}{
		"id":           e.Id,
		"name":         e.Name,
		"hostname":     e.Hostname,
		"port":         e.Port,
		"username":     e.Username,
		"auth_method":  e.AuthMethod,
		"password_enc": e.PasswordEnc,
		"key_id":       e.KeyId,
		"group_id":     e.GroupId,
		"color":        e.Color,
	}
	if e.CreatedAt != nil {
		m["created_at"] = e.CreatedAt.String()
	}
	if e.UpdatedAt != nil {
		m["updated_at"] = e.UpdatedAt.String()
	}
	return m
}

func keyEntityToMap(e *entity.YdstermKeys) map[string]interface{} {
	if e == nil {
		return nil
	}
	m := map[string]interface{}{
		"id":              e.Id,
		"name":            e.Name,
		"private_key_enc": e.PrivateKeyEnc,
		"public_key":      e.PublicKey,
		"passphrase_enc":  e.PassphraseEnc,
	}
	if e.CreatedAt != nil {
		m["created_at"] = e.CreatedAt.String()
	}
	if e.UpdatedAt != nil {
		m["updated_at"] = e.UpdatedAt.String()
	}
	return m
}

func groupEntityToMap(e *entity.YdstermHostGroups) map[string]interface{} {
	if e == nil {
		return nil
	}
	m := map[string]interface{}{
		"id":         e.Id,
		"name":       e.Name,
		"parent_id":  e.ParentId,
		"sort_order": e.SortOrder,
	}
	if e.CreatedAt != nil {
		m["created_at"] = e.CreatedAt.String()
	}
	if e.UpdatedAt != nil {
		m["updated_at"] = e.UpdatedAt.String()
	}
	return m
}
