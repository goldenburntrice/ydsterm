package dbcore

import (
	"context"
	"sync"
)

type ctxKey string

const syncUserKey ctxKey = "sync_user"

var (
	activeUser     = "LOCALUSER"
	activeUserLock sync.RWMutex
)

func SetActiveUser(username string) {
	activeUserLock.Lock()
	defer activeUserLock.Unlock()
	activeUser = username
}

func GetActiveUser() string {
	activeUserLock.RLock()
	defer activeUserLock.RUnlock()
	return activeUser
}

func GetSyncUser(ctx context.Context) string {
	if v, ok := ctx.Value(syncUserKey).(string); ok && v != "" {
		return v
	}
	return GetActiveUser()
}

func SetSyncUser(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, syncUserKey, username)
}

func WithActiveUser(ctx context.Context) context.Context {
	return SetSyncUser(ctx, GetActiveUser())
}
