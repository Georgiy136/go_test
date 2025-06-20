CREATE TABLE IF NOT EXISTS projects (
    project_id INTEGER PRIMARY KEY,
    name character varying NOT NULL UNIQUE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS goods (
    good_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,
    name character varying NOT NULL,
    description character varying,
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_goods_good_id_project_id PRIMARY KEY (good_id, project_id)
);

CREATE INDEX IF NOT EXISTS idx_goods_project_id on goods (project_id);

CREATE INDEX IF NOT EXISTS idx_goods_good_id_project_id on goods (good_id, project_id);

CREATE SEQUENCE IF NOT EXISTS good_sq AS INTEGER START WITH 1;

CREATE SEQUENCE IF NOT EXISTS project_sq AS INTEGER START WITH 1;