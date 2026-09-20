CREATE TABLE user_total_reviews (
    user_id TEXT NOT NULL PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

    track_count INTEGER NOT NULL DEFAULT 0,
    listening_time INTEGER NOT NULL DEFAULT 0,
    avg_completion INTEGER NOT NULL DEFAULT 0,
    skip_count INTEGER NOT NULL DEFAULT 0,
    unique_tracks INTEGER NOT NULL DEFAULT 0,
    favorite_plays INTEGER NOT NULL DEFAULT 0,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE user_total_review_tracks (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id TEXT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,

    rank INTEGER NOT NULL,
    play_count INTEGER NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    FOREIGN KEY (user_id) REFERENCES user_total_reviews(user_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, track_id)
);

CREATE TABLE user_total_review_albums (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    album_id TEXT NOT NULL REFERENCES albums(id) ON DELETE CASCADE,

    rank INTEGER NOT NULL,
    play_count INTEGER NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    FOREIGN KEY (user_id) REFERENCES user_total_reviews(user_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, album_id)
);

CREATE TABLE user_total_review_artists (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    artist_id TEXT NOT NULL REFERENCES artists(id) ON DELETE CASCADE,

    rank INTEGER NOT NULL,
    play_count INTEGER NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    FOREIGN KEY (user_id) REFERENCES user_total_reviews(user_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, artist_id)
);

CREATE TABLE user_total_review_tags (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tag_slug TEXT NOT NULL REFERENCES tags(slug) ON DELETE CASCADE,

    rank INTEGER NOT NULL,
    play_count INTEGER NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    FOREIGN KEY (user_id) REFERENCES user_total_reviews(user_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, tag_slug)
);

CREATE TABLE user_total_review_decades (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    decade INTEGER NOT NULL,

    rank INTEGER NOT NULL,
    play_count INTEGER NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    FOREIGN KEY (user_id) REFERENCES user_total_reviews(user_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, decade)
);

CREATE TABLE user_total_review_months (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    month INTEGER NOT NULL,

    play_count INTEGER NOT NULL DEFAULT 0,
    play_time INTEGER NOT NULL,

    avg_completion INTEGER NOT NULL DEFAULT 0,
    skip_count INTEGER NOT NULL DEFAULT 0,
    unique_tracks INTEGER NOT NULL DEFAULT 0,
    favorite_plays INTEGER NOT NULL DEFAULT 0,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    FOREIGN KEY (user_id) REFERENCES user_total_reviews(user_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, month)
);