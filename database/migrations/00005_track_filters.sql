CREATE TABLE track_filters (
    id TEXT NOT NULL PRIMARY KEY,

    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    name TEXT NOT NULL CHECK(name<>''),
    filter TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);
