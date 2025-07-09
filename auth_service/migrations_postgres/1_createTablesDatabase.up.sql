CREATE SCHEMA IF NOT EXISTS login

CREATE TABLE IF NOT EXISTS login.user_login (
    user_id integer NOT NULL,
    user_agent character varying NOT NULL,
    refresh_token character varying NOT NULL,
    ip_address character varying NOT NULL,
    sign_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user (
    user_id INTEGER PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    token_id INTEGER PRIMARY KEY,
    refresh_token TEXT
);

CREATE INDEX IF NOT EXISTS idx_user_login_user_id on user_login (user_id);

CREATE SEQUENCE IF NOT EXISTS refresh_tokens_sq AS INTEGER START WITH 1;