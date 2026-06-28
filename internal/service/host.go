package service

import (
	"context"

	"ydsterm/internal/crypto"
	"ydsterm/internal/dbcore/dao"
	"ydsterm/internal/types"
)

type HostServiceImpl struct{}

func NewHostService() *HostServiceImpl { return &HostServiceImpl{} }

func (s *HostServiceImpl) Create(input types.HostCreateInput) (*types.Host, error) {
	ctx := context.Background()
	passwordEnc := ""
	if input.Password != "" {
		var err error
		passwordEnc, err = crypto.Encrypt(input.Password)
		if err != nil {
			return nil, err
		}
	}
	input.Password = "" // clear plaintext
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
	return entityToHost(host), nil
}

func (s *HostServiceImpl) Update(input types.HostUpdateInput) (*types.Host, error) {
	ctx := context.Background()
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
	return s.Get(input.ID)
}

func (s *HostServiceImpl) Delete(id string) error {
	return dao.YdstermHosts.Delete(context.Background(), id)
}

func (s *HostServiceImpl) Get(id string) (*types.Host, error) {
	host, err := dao.YdstermHosts.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return entityToHost(host), nil
}

func (s *HostServiceImpl) List(groupID string) ([]types.Host, error) {
	hosts, err := dao.YdstermHosts.List(context.Background(), groupID)
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
	ctx := context.Background()
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
	return entityToKey(key), nil
}

func (s *HostServiceImpl) UpdateKey(input types.KeyUpdateInput) (*types.Key, error) {
	ctx := context.Background()
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
	return s.GetKey(input.ID)
}

func (s *HostServiceImpl) DeleteKey(id string) error {
	return dao.YdstermKeys.Delete(context.Background(), id)
}

func (s *HostServiceImpl) GetKey(id string) (*types.Key, error) {
	key, err := dao.YdstermKeys.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return entityToKey(key), nil
}

func (s *HostServiceImpl) ListKeys() ([]types.Key, error) {
	keys, err := dao.YdstermKeys.List(context.Background())
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
	g, err := dao.YdstermHostGroups.Create(context.Background(), input)
	if err != nil {
		return nil, err
	}
	return entityToHostGroup(g), nil
}

func (s *HostServiceImpl) UpdateGroup(input types.HostGroupUpdateInput) (*types.HostGroup, error) {
	ctx := context.Background()
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
	g, err := dao.YdstermHostGroups.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return entityToHostGroup(g), nil
}

func (s *HostServiceImpl) DeleteGroup(id string) error {
	return dao.YdstermHostGroups.Delete(context.Background(), id)
}

func (s *HostServiceImpl) ListGroups() ([]types.HostGroup, error) {
	groups, err := dao.YdstermHostGroups.List(context.Background())
	if err != nil {
		return nil, err
	}
	out := make([]types.HostGroup, len(groups))
	for i, g := range groups {
		out[i] = *entityToHostGroup(&g)
	}
	return out, nil
}
