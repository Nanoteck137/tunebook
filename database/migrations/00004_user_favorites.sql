-- +goose Up
CREATE TABLE user_favorites (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id TEXT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,

    added INTEGER NOT NULL,

    PRIMARY KEY(user_id, track_id)
);
