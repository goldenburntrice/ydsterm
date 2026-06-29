package service

import (
	"context"
	"fmt"
	"log"

	"ydsterm/internal/dbcore"
	"ydsterm/internal/dbcore/dao"
	"ydsterm/internal/types"
)

type UserServiceImpl struct{}

func NewUserService() *UserServiceImpl { return &UserServiceImpl{} }

func (s *UserServiceImpl) GetCurrentUser() (*types.SyncUser, error) {
	ctx := context.Background()
	u, err := dao.YdstermSyncUsers.GetActive(ctx)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("no active user")
	}
	return syncUserEntityToTypes(u), nil
}

func (s *UserServiceImpl) ListUsers() ([]types.SyncUser, error) {
	ctx := context.Background()
	users, err := dao.YdstermSyncUsers.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]types.SyncUser, len(users))
	for i, u := range users {
		out[i] = *syncUserEntityToTypes(&u)
	}
	return out, nil
}

func (s *UserServiceImpl) SwitchUser(userID string) error {
	ctx := context.Background()
	u, err := dao.YdstermSyncUsers.GetActive(ctx)
	if err != nil {
		return err
	}
	if u != nil && u.Id == userID {
		return nil
	}
	if err := dao.YdstermSyncUsers.SetActive(ctx, userID); err != nil {
		return err
	}
	newUser, err := dao.YdstermSyncUsers.GetActive(ctx)
	if err != nil {
		return err
	}
	if newUser != nil {
		dbcore.SetActiveUser(newUser.Username)
		log.Printf("switched to user: %s", newUser.Username)
	}
	if syncService != nil && newUser != nil && newUser.Username != "LOCALUSER" {
		go func() {
			if err := syncService.PullAll(); err != nil {
				log.Printf("sync pull after switch: %v", err)
			}
		}()
	}
	return nil
}

func (s *UserServiceImpl) RegisterVerifiedUser(username string, passwordHash string) (*types.SyncUser, error) {
	ctx := context.Background()
	existing, err := dao.YdstermSyncUsers.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		if err := dao.YdstermSyncUsers.SetActive(ctx, existing.Id); err != nil {
			return nil, err
		}
		if passwordHash != "" {
			_ = dao.YdstermSyncUsers.UpdatePasswordHash(ctx, existing.Id, passwordHash)
		}
		dbcore.SetActiveUser(username)
		if syncService != nil && username != "LOCALUSER" {
			go func() {
				if err := syncService.PullAll(); err != nil {
					log.Printf("sync pull after verify: %v", err)
				}
			}()
		}
		return syncUserEntityToTypes(existing), nil
	}
	u, err := dao.YdstermSyncUsers.Create(ctx, username, passwordHash)
	if err != nil {
		return nil, err
	}
	if err := dao.YdstermSyncUsers.SetActive(ctx, u.Id); err != nil {
		return nil, err
	}
	dbcore.SetActiveUser(username)
	if syncService != nil && username != "LOCALUSER" {
		go func() {
			if err := syncService.PullAll(); err != nil {
				log.Printf("sync pull after verify: %v", err)
			}
		}()
	}
	return syncUserEntityToTypes(u), nil
}

func (s *UserServiceImpl) InitActiveUser() error {
	ctx := context.Background()
	u, err := dao.YdstermSyncUsers.GetActive(ctx)
	if err != nil {
		return err
	}
	if u != nil {
		dbcore.SetActiveUser(u.Username)
		return nil
	}
	users, err := dao.YdstermSyncUsers.List(ctx)
	if err != nil {
		return err
	}
	if len(users) > 0 {
		if err := dao.YdstermSyncUsers.SetActive(ctx, users[0].Id); err != nil {
			return err
		}
		dbcore.SetActiveUser(users[0].Username)
		return nil
	}
	newUser, err := dao.YdstermSyncUsers.Create(ctx, "LOCALUSER", "")
	if err != nil {
		return err
	}
	if err := dao.YdstermSyncUsers.SetActive(ctx, newUser.Id); err != nil {
		return err
	}
	dbcore.SetActiveUser("LOCALUSER")
	return nil
}

func syncUserEntityToTypes(e *dao.SyncUserEntity) *types.SyncUser {
	if e == nil {
		return nil
	}
	u := &types.SyncUser{
		ID:       e.Id,
		Username: e.Username,
		IsActive: e.IsActive == 1,
		CreatedAt: gtimeToTime(e.CreatedAt),
		UpdatedAt: gtimeToTime(e.UpdatedAt),
	}
	if e.LastSyncAt != nil {
		t := gtimeToTime(e.LastSyncAt)
		u.LastSyncAt = &t
	}
	return u
}
