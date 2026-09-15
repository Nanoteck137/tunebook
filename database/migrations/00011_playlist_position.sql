ALTER TABLE playlists ADD COLUMN position INTEGER NOT NULL DEFAULT 0;

UPDATE playlists
SET position = (
    SELECT COUNT(*)
    FROM playlists AS p2
    WHERE p2.created < playlists.created
);
