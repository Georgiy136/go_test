CREATE TABLE IF NOT EXISTS user_login (
    user_id integer NOT NULL,
    user_agent character varying NOT NULL,
    refresh_token character varying NOT NULL,
    ip_address character varying NOT NULL,
    sign_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user (
    user_id INTEGER PRIMARY KEY
);

CREATE INDEX IF NOT EXISTS idx_user_login_user_id on user_login (user_id);