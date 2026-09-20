-- Indexes derived from EXPLAIN QUERY PLAN analysis of the app's hot query paths.
--
-- user_track_stats: the PK is (user_id, track_id, period_type, year, period_value),
-- so period-filtered aggregation queries (all/year/month summaries plus every
-- review-generation query) had to scan the whole user's rows across all four
-- period buckets. A (user_id, period_type, ...) leading index lets them jump
-- straight to the relevant bucket.
CREATE INDEX idx_user_track_stats_user_period ON user_track_stats(user_id, period_type, year, period_value);

-- GetUserTopTracks orders by play_count DESC. This index lets the top-tracks
-- query walk the bucket in rank order (and doubles as a covering index for the
-- period-only aggregation queries).
CREATE INDEX idx_user_track_stats_user_period_play ON user_track_stats(user_id, period_type, play_count DESC);

-- User review lists are read with "WHERE user_id = ? ORDER BY rank" (+ year for
-- the year reviews). The PKs are (user_id, object_id), so rank-ordered
-- pagination previously required a temp b-tree sort per page.
CREATE INDEX idx_user_total_review_tracks_rank ON user_total_review_tracks(user_id, rank);
CREATE INDEX idx_user_total_review_albums_rank ON user_total_review_albums(user_id, rank);
CREATE INDEX idx_user_total_review_artists_rank ON user_total_review_artists(user_id, rank);
CREATE INDEX idx_user_total_review_tags_rank ON user_total_review_tags(user_id, rank);
CREATE INDEX idx_user_total_review_decades_rank ON user_total_review_decades(user_id, rank);

CREATE INDEX idx_user_year_review_tracks_rank ON user_year_review_tracks(user_id, year, rank);
CREATE INDEX idx_user_year_review_albums_rank ON user_year_review_albums(user_id, year, rank);
CREATE INDEX idx_user_year_review_artists_rank ON user_year_review_artists(user_id, year, rank);
CREATE INDEX idx_user_year_review_tags_rank ON user_year_review_tags(user_id, year, rank);
CREATE INDEX idx_user_year_review_decades_rank ON user_year_review_decades(user_id, year, rank);

-- GetPendingJobs filters by status and pulls the oldest first.
CREATE INDEX idx_jobs_status_created ON jobs(status, created);

-- GetUserPlaylists / reorder helpers order a user's playlists by position.
CREATE INDEX idx_playlists_owner_position ON playlists(owner_id, position);

-- GetNextTrackFilterPosition and the filters list order by position per user.
CREATE INDEX idx_track_filters_user_position ON track_filters(user_id, position);