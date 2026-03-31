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

CREATE TABLE playlist_filters (
    id TEXT NOT NULL,
    playlist_id TEXT NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,

    name TEXT NOT NULL CHECK(name<>''),
    filter TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    PRIMARY KEY(id, playlist_id)
);

CREATE TABLE track_filters (
    id TEXT NOT NULL PRIMARY KEY,

    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    name TEXT NOT NULL CHECK(name<>''),
    filter TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

ALTER TABLE users_settings ADD COLUMN quick_playlist TEXT REFERENCES playlists(id) ON DELETE SET NULL;
