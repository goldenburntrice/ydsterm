CREATE TABLE IF NOT EXISTS ydsterm_schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ydsterm_host_groups (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    parent_id TEXT DEFAULT '',
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ydsterm_keys (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    private_key_enc TEXT NOT NULL,
    public_key TEXT DEFAULT '',
    passphrase_enc TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ydsterm_hosts (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    hostname TEXT NOT NULL,
    port INTEGER DEFAULT 22,
    username TEXT DEFAULT '',
    auth_method TEXT DEFAULT 'password',
    password_enc TEXT DEFAULT '',
    key_id TEXT DEFAULT '',
    group_id TEXT DEFAULT '',
    color TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

CREATE TABLE IF NOT EXISTS ydsterm_snippet_folders (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    parent_id TEXT DEFAULT '',
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ydsterm_snippets (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    content TEXT NOT NULL,
    language TEXT DEFAULT '',
    folder_id TEXT DEFAULT '',
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ydsterm_port_forwards (
    id TEXT PRIMARY KEY,
    name TEXT DEFAULT '',
    host_id TEXT NOT NULL,
    type TEXT NOT NULL,
    local_address TEXT DEFAULT '127.0.0.1',
    local_port INTEGER DEFAULT 0,
    remote_host TEXT DEFAULT '',
    remote_port INTEGER DEFAULT 0,
    socks_host TEXT DEFAULT '127.0.0.1',
    socks_port INTEGER DEFAULT 0,
    enabled INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
