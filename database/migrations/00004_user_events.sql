-- +goose Up
CREATE TABLE user_listening_events (
    id          TEXT PRIMARY KEY,

    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id    TEXT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,

    listened_at INTEGER NOT NULL,
    percent REAL NOT NULL,
    position_ms INTEGER NOT NULL, 
    source TEXT NOT NULL
);

-- CREATE INDEX idx_listen_events_user_time ON listen_events(user_id, listened_at);
-- CREATE INDEX idx_listen_events_track ON listen_events(track_id);
