ALTER TABLE tracks ADD COLUMN album_position INTEGER NOT NULL DEFAULT 0;

-- Serves the album queries, which filter on album_id and order by
-- album_position. Mirrors the old idx_tracks_album_number(album_id, number).
CREATE INDEX idx_tracks_album_position ON tracks(album_id, album_position);
