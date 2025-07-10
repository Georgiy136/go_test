CREATE SCHEMA IF NOT EXISTS users;

CREATE SCHEMA IF NOT EXISTS login;

CREATE TABLE IF NOT EXISTS login.user_login (
    user_id integer NOT NULL,
    refresh_token_id character varying NOT NULL,
    user_agent character varying NOT NULL,
    ip_address character varying NOT NULL,
    sign_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users.user (
    user_id INTEGER PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS login.refresh_tokens (
    refresh_token_id INTEGER PRIMARY KEY,
    refresh_token TEXT
);

CREATE SEQUENCE IF NOT EXISTS login.refresh_tokens_sq AS INTEGER START WITH 1;

CREATE INDEX IF NOT EXISTS idx_users_user_id on users.user (user_id);

CREATE INDEX IF NOT EXISTS idx_user_login_user_id on login.user_login (user_id);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_refresh_token_id on login.refresh_tokens (refresh_token_id);
