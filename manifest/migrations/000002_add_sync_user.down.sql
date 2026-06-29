DROP TABLE IF EXISTS ydsterm_settings;
DROP TABLE IF EXISTS ydsterm_sync_users;
ALTER TABLE ydsterm_port_forwards DROP COLUMN sync_user;
ALTER TABLE ydsterm_snippet_folders DROP COLUMN sync_user;
ALTER TABLE ydsterm_snippets DROP COLUMN sync_user;
ALTER TABLE ydsterm_keys DROP COLUMN sync_user;
ALTER TABLE ydsterm_host_groups DROP COLUMN sync_user;
ALTER TABLE ydsterm_hosts DROP COLUMN sync_user;
