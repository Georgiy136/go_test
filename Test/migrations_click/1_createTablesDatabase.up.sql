CREATE TABLE IF NOT EXISTS logs (
    ID SERIAL NOT NULL,
    ProjectID INTEGER NOT NULL,
    Name character varying NOT NULL,
    Description character varying,
    Priority INTEGER,
    Removed BOOLEAN DEFAULT FALSE,
    EventTime TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);