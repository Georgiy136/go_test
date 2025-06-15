CREATE TABLE IF NOT EXISTS projects (
    ID SERIAL PRIMARY KEY,
    name character varying NOT NULL UNIQUE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS goods (
    ID SERIAL NOT NULL,
    project_id INTEGER NOT NULL,
    name character varying NOT NULL,
    description character varying,
    priority INTEGER,
    removed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    CONSTRAINT pk_goods PRIMARY KEY (ID, project_id),
    FOREIGN KEY (project_id) REFERENCES projects (ID)
);

CREATE INDEX idx_goods_project_id on goods (project_id);

CREATE INDEX idx_goods_id_project_id on goods (ID, project_id);