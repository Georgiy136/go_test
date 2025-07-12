CREATE SCHEMA IF NOT EXISTS sessions;

CREATE TABLE IF NOT EXISTS sessions.user_login (
    user_id integer NOT NULL,
    session_id character varying NOT NULL,
    hash_refresh_token character varying NOT NULL,
    user_agent character varying NOT NULL,
    ip_address character varying NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_user_login PRIMARY KEY (user_id, user_agent, ip_address)
);

CREATE INDEX IF NOT EXISTS idx_user_login_user_id on sessions.user_login (user_id);