ALTER TABLE track_filters ADD COLUMN position INTEGER NOT NULL DEFAULT 0;

UPDATE track_filters
SET position = (
    SELECT COUNT(*)
    FROM track_filters AS t2
    WHERE t2.user_id = track_filters.user_id
      AND t2.created < track_filters.created
);