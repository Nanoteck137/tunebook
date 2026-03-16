-- +goose Up
CREATE TABLE artists (
    id TEXT PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,

    name TEXT NOT NULL CHECK(name<>''),
    other_name TEXT,

    cover_art TEXT,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE albums (
    id TEXT PRIMARY KEY,

    name TEXT NOT NULL CHECK(name<>''),
    other_name TEXT,

    artist_id TEXT NOT NULL REFERENCES artists(id),

    cover_art TEXT,
    year INT,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE albums_featuring_artists (
    album_id TEXT NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
    artist_id TEXT NOT NULL REFERENCES artists(id) ON DELETE CASCADE,

    PRIMARY KEY(album_id, artist_id)
);

CREATE TABLE tracks (
    id TEXT PRIMARY KEY,

    filename TEXT NOT NULL,
    modified_time INT NOT NULL,
    media_format TEXT NOT NULL,

    name TEXT NOT NULL CHECK(name<>''),
    other_name TEXT,

    album_id TEXT NOT NULL REFERENCES albums(id),
    artist_id TEXT NOT NULL REFERENCES artists(id),

    duration INT NOT NULL,
    number INT,
    year INT,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE tracks_featuring_artists (
    track_id TEXT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    artist_id TEXT NOT NULL REFERENCES artists(id) ON DELETE CASCADE,

    PRIMARY KEY(track_id, artist_id)
);

CREATE TABLE tags (
    slug TEXT PRIMARY KEY
);

CREATE TABLE artists_tags  (
    artist_id TEXT REFERENCES artists(id) ON DELETE CASCADE,
    tag_slug TEXT REFERENCES tags(slug) ON DELETE CASCADE,

    PRIMARY KEY(artist_id, tag_slug)
);

CREATE TABLE albums_tags  (
    album_id TEXT REFERENCES albums(id) ON DELETE CASCADE,
    tag_slug TEXT REFERENCES tags(slug) ON DELETE CASCADE,

    PRIMARY KEY(album_id, tag_slug)
);

CREATE TABLE tracks_tags (
    track_id TEXT REFERENCES tracks(id) ON DELETE CASCADE,
    tag_slug TEXT REFERENCES tags(slug) ON DELETE CASCADE,

    PRIMARY KEY(track_id, tag_slug)
);

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

    order_num INTEGER NOT NULL,

    -- TODO(patrik): Add created, updated / added FIELD

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
    id TEXT NOT NULL,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    name TEXT NOT NULL CHECK(name<>''),
    filter TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    PRIMARY KEY(id, user_id)
);

CREATE TABLE virtual_playlists (
    id TEXT PRIMARY KEY,

    name TEXT NOT NULL CHECK(name<>''),

    owner_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    playlist_id TEXT REFERENCES playlists(id) ON DELETE CASCADE,

    filter TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE users_settings (
    id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    display_name TEXT,

    quick_playlist TEXT REFERENCES playlists(id) ON DELETE SET NULL
);

CREATE TABLE user_listening_events (
    -- id          INTEGER PRIMARY KEY AUTOINCREMENT,
    id          TEXT PRIMARY KEY,

    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id    TEXT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,

    -- listened_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    listened_at INTEGER NOT NULL,
    -- TODO(patrik): rename to position
    -- TODO(patrik): add percentage of played
    duration_ms INTEGER NOT NULL, -- how long they actually listened
    source      TEXT NOT NULL     -- 'search', 'playlist', 'radio', etc.
);

-- CREATE INDEX idx_listen_events_user_time ON listen_events(user_id, listened_at);
-- CREATE INDEX idx_listen_events_track ON listen_events(track_id);

CREATE TABLE user_listening_stats (
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id    INTEGER NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    period      TEXT NOT NULL,  -- '2024', '2024-Q1', '2024-03'

    play_count  INTEGER NOT NULL DEFAULT 0,
    skip_count  INTEGER NOT NULL DEFAULT 0,
    total_ms    INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY (user_id, track_id, period)
);

-- +goose Down
DROP TABLE api_tokens; 

DROP TABLE taglists; 
DROP TABLE playlist_items; 
DROP TABLE playlists; 

DROP TABLE users_settings; 
DROP TABLE users; 

DROP TABLE tracks_tags; 
DROP TABLE albums_tags ; 
DROP TABLE artists_tags ; 
DROP TABLE tags; 

DROP TABLE tracks_featuring_artists; 
DROP TABLE tracks_media; 
DROP TABLE tracks; 

DROP TABLE albums_featuring_artists; 
DROP TABLE albums; 

DROP TABLE artists; 

