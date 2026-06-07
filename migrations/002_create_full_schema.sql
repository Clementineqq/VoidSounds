CREATE TABLE IF NOT EXISTS users (
    id              SERIAL PRIMARY KEY,
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    name            VARCHAR(100),
    role            VARCHAR(20) NOT NULL DEFAULT 'viewer',
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cities (
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(100) NOT NULL,
    slug    VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS genres (
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(100) UNIQUE NOT NULL,
    slug    VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    date            TIMESTAMP WITH TIME ZONE NOT NULL,
    city_id         INTEGER REFERENCES cities(id),
    address         VARCHAR(255),
    price           INTEGER NOT NULL DEFAULT 0,
    available       INTEGER NOT NULL DEFAULT 0,
    poster_url      VARCHAR(500),
    organizer_id    INTEGER REFERENCES users(id),
    status          VARCHAR(20) NOT NULL DEFAULT 'published',
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS event_genres (
    event_id    INTEGER REFERENCES events(id) ON DELETE CASCADE,
    genre_id    INTEGER REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (event_id, genre_id)
);

CREATE TABLE IF NOT EXISTS tickets (
    id              SERIAL PRIMARY KEY,
    event_id        INTEGER REFERENCES events(id) ON DELETE CASCADE,
    user_id         INTEGER REFERENCES users(id),
    quantity        INTEGER NOT NULL,
    total_price     INTEGER NOT NULL,
    purchase_date   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status          VARCHAR(20) DEFAULT 'paid'
);

CREATE INDEX IF NOT EXISTS idx_events_date ON events(date);
CREATE INDEX IF NOT EXISTS idx_events_city_id ON events(city_id);
CREATE INDEX IF NOT EXISTS idx_events_status ON events(status);



ALTER TABLE users ADD COLUMN is_banned BOOLEAN DEFAULT FALSE;