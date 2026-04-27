-- 001_create_events_table.sql

CREATE TABLE IF NOT EXISTS events (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    date            TIMESTAMP WITH TIME ZONE NOT NULL,
    location        VARCHAR(255) NOT NULL,
    genre           VARCHAR(100),
    price           INTEGER NOT NULL DEFAULT 0,
    available       INTEGER NOT NULL DEFAULT 0,
    organizer_id    INTEGER,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_events_date ON events(date);
CREATE INDEX IF NOT EXISTS idx_events_genre ON events(genre);