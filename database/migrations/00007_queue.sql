-- +goose Up
CREATE TABLE queues (
    id TEXT PRIMARY KEY,
    current_index INTEGER NOT NULL DEFAULT 0,
    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE queue_items (
    id TEXT PRIMARY KEY,
    queue_id TEXT NOT NULL REFERENCES queues(id) ON DELETE CASCADE,
    track_id TEXT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    created INTEGER NOT NULL
);
