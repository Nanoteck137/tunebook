CREATE TABLE user_stats (
    user_id TEXT NOT NULL PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

    num_tracks_played INTEGER NOT NULL DEFAULT 0,
    num_tracks_skipped INTEGER NOT NULL DEFAULT 0,
    num_playlists_created INTEGER NOT NULL DEFAULT 0,
    num_favorite_tracks INTEGER NOT NULL DEFAULT 0,
    listening_time INTEGER NOT NULL DEFAULT 0,
    last_listened_at INTEGER,

    updated INTEGER NOT NULL
);
