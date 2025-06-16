CREATE TABLE IF NOT EXISTS logs (
    LogID SERIAL NOT NULL,
    GoodsID INTEGER NOT NULL,
    ProjectID INTEGER NOT NULL,
    Name character varying NOT NULL,
    Description character varying,
    Priority INTEGER DEFAULT 0,
    Removed BOOLEAN DEFAULT FALSE,
    EventTime TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);