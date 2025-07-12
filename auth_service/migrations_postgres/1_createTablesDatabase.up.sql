CREATE SCHEMA IF NOT EXISTS login;

CREATE TABLE IF NOT EXISTS login.user_login (
    user_id integer NOT NULL,
    refresh_token_id character varying NOT NULL,
    user_agent character varying NOT NULL,
    ip_address character varying NOT NULL,
    CONSTRAINT pk_user_login PRIMARY KEY (user_id, user_agent, ip_address)
);

CREATE TABLE IF NOT EXISTS login.refresh_tokens (
    refresh_token_id SERIAL PRIMARY KEY,
    refresh_token TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_login_user_id on login.user_login (user_id);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_refresh_token_id on login.refresh_tokens (refresh_token_id);
