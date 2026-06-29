package service

import (
	"ydsterm/internal/dbcore/dao"
	"ydsterm/internal/dbcore/model/entity"
	"ydsterm/internal/types"
)

type PortForwardServiceImpl struct{}

func NewPortForwardService() *PortForwardServiceImpl { return &PortForwardServiceImpl{} }

func (s *PortForwardServiceImpl) Create(input types.PortForwardCreateInput) (*types.PortForward, error) {
	pf, err := dao.YdstermPortForwards.Create(activeCtx(), input)
	if err != nil {
		return nil, err
	}
	if syncService != nil {
		go syncService.PushPortForward(portForwardEntityToMap(pf))
	}
	return entityToPortForward(pf), nil
}

func (s *PortForwardServiceImpl) Update(input types.PortForwardUpdateInput) (*types.PortForward, error) {
	ctx := activeCtx()
	data := make(map[string]interface{})
	cols := dao.YdstermPortForwards.Columns()

	if input.Name != nil {
		data[cols.Name] = *input.Name
	}
	if input.HostID != nil {
		data[cols.HostId] = *input.HostID
	}
	if input.Type != nil {
		data[cols.Type] = *input.Type
	}
	if input.LocalAddress != nil {
		data[cols.LocalAddress] = *input.LocalAddress
	}
	if input.LocalPort != nil {
		data[cols.LocalPort] = *input.LocalPort
	}
	if input.RemoteHost != nil {
		data[cols.RemoteHost] = *input.RemoteHost
	}
	if input.RemotePort != nil {
		data[cols.RemotePort] = *input.RemotePort
	}
	if input.SocksHost != nil {
		data[cols.SocksHost] = *input.SocksHost
	}
	if input.SocksPort != nil {
		data[cols.SocksPort] = *input.SocksPort
	}

	if err := dao.YdstermPortForwards.Update(ctx, input.ID, data); err != nil {
		return nil, err
	}
	pf, _ := dao.YdstermPortForwards.Get(ctx, input.ID)
	if syncService != nil && pf != nil {
		go syncService.PushPortForward(portForwardEntityToMap(pf))
	}
	return entityToPortForward(pf), nil
}

func (s *PortForwardServiceImpl) Delete(id string) error {
	if err := dao.YdstermPortForwards.Delete(activeCtx(), id); err != nil {
		return err
	}
	if syncService != nil {
		go syncService.DeletePortForward(id)
	}
	return nil
}

func (s *PortForwardServiceImpl) Get(id string) (*types.PortForward, error) {
	pf, err := dao.YdstermPortForwards.Get(activeCtx(), id)
	if err != nil {
		return nil, err
	}
	return entityToPortForward(pf), nil
}

func (s *PortForwardServiceImpl) List(hostID string) ([]types.PortForward, error) {
	pfs, err := dao.YdstermPortForwards.List(activeCtx(), hostID)
	if err != nil {
		return nil, err
	}
	out := make([]types.PortForward, len(pfs))
	for i, p := range pfs {
		out[i] = *entityToPortForward(&p)
	}
	return out, nil
}

func (s *PortForwardServiceImpl) Toggle(id string, enabled bool) error {
	if err := dao.YdstermPortForwards.Toggle(activeCtx(), id, enabled); err != nil {
		return err
	}
	if syncService != nil {
		pf, _ := dao.YdstermPortForwards.Get(activeCtx(), id)
		if pf != nil {
			go syncService.PushPortForward(portForwardEntityToMap(pf))
		}
	}
	return nil
}

func portForwardEntityToMap(e *entity.YdstermPortForwards) map[string]interface{} {
	if e == nil {
		return nil
	}
	m := map[string]interface{}{
		"id":            e.Id,
		"name":          e.Name,
		"host_id":       e.HostId,
		"type":          e.Type,
		"local_address": e.LocalAddress,
		"local_port":    e.LocalPort,
		"remote_host":   e.RemoteHost,
		"remote_port":   e.RemotePort,
		"socks_host":    e.SocksHost,
		"socks_port":    e.SocksPort,
		"enabled":       e.Enabled,
	}
	if e.CreatedAt != nil {
		m["created_at"] = e.CreatedAt.String()
	}
	if e.UpdatedAt != nil {
		m["updated_at"] = e.UpdatedAt.String()
	}
	return m
}
