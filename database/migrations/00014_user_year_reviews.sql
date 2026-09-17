CREATE TABLE user_year_reviews (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,

    track_count INTEGER NOT NULL DEFAULT 0,
    listening_time INTEGER NOT NULL DEFAULT 0,
    avg_completion INTEGER NOT NULL DEFAULT 0,
    skip_count INTEGER NOT NULL DEFAULT 0,
    unique_tracks INTEGER NOT NULL DEFAULT 0,
    favorite_plays INTEGER NOT NULL DEFAULT 0,

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    PRIMARY KEY (user_id, year)
);

CREATE TABLE user_year_review_tracks (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,
    track_id TEXT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,

    rank INTEGER NOT NULL,
    play_count INTEGER NOT NULL,

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    FOREIGN KEY (user_id, year) REFERENCES user_year_reviews(user_id, year) ON DELETE CASCADE,
    PRIMARY KEY (user_id, year, track_id)
);

CREATE TABLE user_year_review_albums (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,
    album_id TEXT NOT NULL REFERENCES albums(id) ON DELETE CASCADE,

    rank INTEGER NOT NULL,
    play_count INTEGER NOT NULL,

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    FOREIGN KEY (user_id, year) REFERENCES user_year_reviews(user_id, year) ON DELETE CASCADE,
    PRIMARY KEY (user_id, year, album_id)
);

CREATE TABLE user_year_review_artists (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,
    artist_id TEXT NOT NULL REFERENCES artists(id) ON DELETE CASCADE,

    rank INTEGER NOT NULL,
    play_count INTEGER NOT NULL,

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    FOREIGN KEY (user_id, year) REFERENCES user_year_reviews(user_id, year) ON DELETE CASCADE,
    PRIMARY KEY (user_id, year, artist_id)
);

CREATE TABLE user_year_review_tags (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,
    tag_slug TEXT NOT NULL REFERENCES tags(slug) ON DELETE CASCADE,

    rank INTEGER NOT NULL,
    play_count INTEGER NOT NULL,

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    FOREIGN KEY (user_id, year) REFERENCES user_year_reviews(user_id, year) ON DELETE CASCADE,
    PRIMARY KEY (user_id, year, tag_slug)
);

CREATE TABLE user_year_review_decades (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,
    decade INTEGER NOT NULL,

    rank INTEGER NOT NULL,
    play_count INTEGER NOT NULL,

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    FOREIGN KEY (user_id, year) REFERENCES user_year_reviews(user_id, year) ON DELETE CASCADE,
    PRIMARY KEY (user_id, year, decade)
);

CREATE TABLE user_year_review_months (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,
    month INTEGER NOT NULL,

    play_count INTEGER NOT NULL DEFAULT 0,
    play_time INTEGER NOT NULL,
    avg_completion INTEGER NOT NULL DEFAULT 0,
    skip_count INTEGER NOT NULL DEFAULT 0,
    unique_tracks INTEGER NOT NULL DEFAULT 0,
    favorite_plays INTEGER NOT NULL DEFAULT 0,

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    FOREIGN KEY (user_id, year) REFERENCES user_year_reviews(user_id, year) ON DELETE CASCADE,
    PRIMARY KEY (user_id, year, month)
);

CREATE TABLE user_year_review_month_tracks (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,
    month INTEGER NOT NULL,
    track_id TEXT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,

    rank INTEGER NOT NULL,
    play_count INTEGER NOT NULL,

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    FOREIGN KEY (user_id, year) REFERENCES user_year_reviews(user_id, year) ON DELETE CASCADE,
    PRIMARY KEY (user_id, year, month, track_id)
);

CREATE TABLE user_year_review_month_albums (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,
    month INTEGER NOT NULL,
    album_id TEXT NOT NULL REFERENCES albums(id) ON DELETE CASCADE,

    rank INTEGER NOT NULL,
    play_count INTEGER NOT NULL,

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    FOREIGN KEY (user_id, year) REFERENCES user_year_reviews(user_id, year) ON DELETE CASCADE,
    PRIMARY KEY (user_id, year, month, album_id)
);

CREATE TABLE user_year_review_month_artists (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,
    month INTEGER NOT NULL,
    artist_id TEXT NOT NULL REFERENCES artists(id) ON DELETE CASCADE,

    rank INTEGER NOT NULL,
    play_count INTEGER NOT NULL,

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    FOREIGN KEY (user_id, year) REFERENCES user_year_reviews(user_id, year) ON DELETE CASCADE,
    PRIMARY KEY (user_id, year, month, artist_id)
);

CREATE TABLE user_year_review_month_tags (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,
    month INTEGER NOT NULL,
    tag_slug TEXT NOT NULL REFERENCES tags(slug) ON DELETE CASCADE,

    rank INTEGER NOT NULL,
    play_count INTEGER NOT NULL,

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    FOREIGN KEY (user_id, year) REFERENCES user_year_reviews(user_id, year) ON DELETE CASCADE,
    PRIMARY KEY (user_id, year, month, tag_slug)
);

CREATE TABLE user_year_review_month_decades (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,
    month INTEGER NOT NULL,
    decade INTEGER NOT NULL,

    rank INTEGER NOT NULL,
    play_count INTEGER NOT NULL,

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    FOREIGN KEY (user_id, year) REFERENCES user_year_reviews(user_id, year) ON DELETE CASCADE,
    PRIMARY KEY (user_id, year, month, decade)
);
