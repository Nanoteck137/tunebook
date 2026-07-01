CREATE TABLE queues (
    id TEXT NOT NULL,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    current_index INTEGER NOT NULL DEFAULT 0,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    PRIMARY KEY (id, user_id)
);

CREATE TABLE queue_items (
    id TEXT PRIMARY KEY,

    queue_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    track_id TEXT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,

    position INTEGER NOT NULL,

    created INTEGER NOT NULL,

    FOREIGN KEY (queue_id, user_id) REFERENCES queues(id, user_id) ON DELETE CASCADE
);
