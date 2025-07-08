CREATE TABLE IF NOT EXISTS user_login (
    user_id integer NOT NULL,
    user_agent character varying NOT NULL,
    refresh_token character varying NOT NULL,
    ip_address character varying NOT NULL,
    sign_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS goods (
    good_id INTEGER,
    project_id INTEGER,
    name character varying,
    description character varying,
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (good_id, project_id)
);

CREATE INDEX IF NOT EXISTS idx_goods_project_id on goods (project_id);

CREATE INDEX IF NOT EXISTS idx_goods_good_id_project_id on goods (good_id, project_id);

CREATE SEQUENCE IF NOT EXISTS good_sq AS INTEGER START WITH 1;