CREATE INDEX IF NOT EXISTS idx_events_city_id ON events(city_id);
CREATE INDEX IF NOT EXISTS idx_events_date ON events(date);
CREATE INDEX IF NOT EXISTS idx_events_status ON events(status);

CREATE INDEX IF NOT EXISTS idx_event_genres_event_id ON event_genres(event_id);
CREATE INDEX IF NOT EXISTS idx_event_genres_genre_id ON event_genres(genre_id);

CREATE INDEX IF NOT EXISTS idx_tickets_user_id ON tickets(user_id);