-- +goose Up
CREATE TABLE playlists (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(name<>''),
    cover_art TEXT,

    owner_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE playlist_items (
    playlist_id TEXT NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    track_id TEXT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,

    position INTEGER NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    PRIMARY KEY(playlist_id, track_id)
);

ALTER TABLE users_settings ADD COLUMN quick_playlist TEXT REFERENCES playlists(id) ON DELETE SET NULL;
