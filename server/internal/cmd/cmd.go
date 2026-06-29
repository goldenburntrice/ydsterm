package cmd

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"ydsterm-server/internal/consts"
	"ydsterm-server/internal/controller"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "YDSterm sync server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			if err := initDatabase(ctx); err != nil {
				log.Fatalf("init database: %v", err)
			}

			serverKey, err := ensureServerKey(ctx)
			if err != nil {
				log.Fatalf("ensure server key: %v", err)
			}
			log.Printf("Server key: %s", serverKey)

			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(serverKeyAuthMiddleware(serverKey))
				group.POST("/api/auth/verify", controller.Auth.Verify)
				group.GET("/api/sync/pull", controller.Sync.Pull)
				group.POST("/api/sync/push", controller.Sync.Push)
			})

			s.Run()
			return nil
		},
	}
)

func initDatabase(ctx context.Context) error {
	db := g.DB()

	if _, err := db.Exec(ctx, `CREATE TABLE IF NOT EXISTS server_config (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, `CREATE TABLE IF NOT EXISTS sync_users (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return err
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS sync_hosts (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL DEFAULT '',
			hostname TEXT NOT NULL DEFAULT '',
			port INTEGER DEFAULT 22,
			username TEXT DEFAULT '',
			auth_method TEXT DEFAULT 'password',
			password_enc TEXT DEFAULT '',
			key_id TEXT DEFAULT '',
			group_id TEXT DEFAULT '',
			color TEXT DEFAULT '',
			sync_user TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sync_host_groups (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL DEFAULT '',
			parent_id TEXT DEFAULT '',
			sort_order INTEGER DEFAULT 0,
			sync_user TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sync_keys (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL DEFAULT '',
			private_key_enc TEXT DEFAULT '',
			public_key TEXT DEFAULT '',
			passphrase_enc TEXT DEFAULT '',
			sync_user TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sync_snippets (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL DEFAULT '',
			content TEXT DEFAULT '',
			language TEXT DEFAULT '',
			folder_id TEXT DEFAULT '',
			sort_order INTEGER DEFAULT 0,
			sync_user TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sync_snippet_folders (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL DEFAULT '',
			parent_id TEXT DEFAULT '',
			sort_order INTEGER DEFAULT 0,
			sync_user TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sync_port_forwards (
			id TEXT PRIMARY KEY,
			name TEXT DEFAULT '',
			host_id TEXT NOT NULL DEFAULT '',
			type TEXT NOT NULL DEFAULT '',
			local_address TEXT DEFAULT '127.0.0.1',
			local_port INTEGER DEFAULT 0,
			remote_host TEXT DEFAULT '',
			remote_port INTEGER DEFAULT 0,
			socks_host TEXT DEFAULT '127.0.0.1',
			socks_port INTEGER DEFAULT 0,
			enabled INTEGER DEFAULT 0,
			sync_user TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, sql := range tables {
		if _, err := db.Exec(ctx, sql); err != nil {
			return err
		}
	}

	log.Println("Database initialized")
	return nil
}

func ensureServerKey(ctx context.Context) (string, error) {
	cfgKey := g.Cfg().MustGet(ctx, "server-config.serverKey", "").String()
	if cfgKey != "" {
		_, err := g.DB().Exec(ctx, "INSERT OR REPLACE INTO server_config (key, value) VALUES ('server_key', ?)", cfgKey)
		return cfgKey, err
	}

	record, err := g.DB().Model("server_config").Ctx(ctx).Where("key", "server_key").One()
	if err != nil {
		return "", err
	}
	if !record.IsEmpty() {
		return record["value"].String(), nil
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	key := hex.EncodeToString(b)

	_, err = g.DB().Exec(ctx, "INSERT INTO server_config (key, value) VALUES ('server_key', ?)", key)
	return key, err
}

func serverKeyAuthMiddleware(expectedKey string) func(r *ghttp.Request) {
	return func(r *ghttp.Request) {
		key := r.Header.Get("X-Server-Key")
		if key == "" {
			r.Response.WriteJsonExit(g.Map{"error": "missing X-Server-Key header"})
			return
		}
		if key != expectedKey {
			r.Response.WriteStatusExit(401, g.Map{"error": "invalid server key"})
			return
		}
		r.Middleware.Next()
	}
}

func init() {
	_ = consts.LocalUser
}
