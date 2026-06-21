CREATE TABLE track_history (
    id TEXT NOT NULL PRIMARY KEY,

    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id TEXT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,

    listened_at INTEGER NOT NULL,

    playback_type TEXT NOT NULL,
    status TEXT NOT NULL,

    percent_played INTEGER NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);
