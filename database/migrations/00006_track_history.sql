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

CREATE TABLE user_track_stats (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id TEXT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,

    -- time bucket
    period_type TEXT NOT NULL, -- all, year, quarter, month
    year INTEGER, -- 2026
    period_value INTEGER, -- all = NULL, year = NULL, quarter = 1-4, month = 1-12

    -- metrics
    play_count INTEGER NOT NULL DEFAULT 0,
    skip_count INTEGER NOT NULL DEFAULT 0,
    play_time INTEGER NOT NULL DEFAULT 0,

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    PRIMARY KEY (user_id, track_id, period_type, year, period_value)
);
