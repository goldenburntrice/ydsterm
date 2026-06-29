CREATE TABLE IF NOT EXISTS server_config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sync_users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sync_hosts (
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
);

CREATE TABLE IF NOT EXISTS sync_host_groups (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',
    parent_id TEXT DEFAULT '',
    sort_order INTEGER DEFAULT 0,
    sync_user TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sync_keys (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',
    private_key_enc TEXT DEFAULT '',
    public_key TEXT DEFAULT '',
    passphrase_enc TEXT DEFAULT '',
    sync_user TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sync_snippets (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',
    content TEXT DEFAULT '',
    language TEXT DEFAULT '',
    folder_id TEXT DEFAULT '',
    sort_order INTEGER DEFAULT 0,
    sync_user TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sync_snippet_folders (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',
    parent_id TEXT DEFAULT '',
    sort_order INTEGER DEFAULT 0,
    sync_user TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sync_port_forwards (
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
);
