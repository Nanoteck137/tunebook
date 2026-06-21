CREATE TABLE artists (
    id TEXT PRIMARY KEY,

    name TEXT NOT NULL CHECK(name<>''),

    cover_art TEXT,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE albums (
    id TEXT PRIMARY KEY,

    name TEXT NOT NULL CHECK(name<>''),

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
