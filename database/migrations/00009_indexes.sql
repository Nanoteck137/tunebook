CREATE INDEX idx_albums_artist_id ON albums(artist_id);
CREATE INDEX idx_tracks_album_id ON tracks(album_id);
CREATE INDEX idx_tracks_artist_id ON tracks(artist_id);
CREATE INDEX idx_playlists_owner_id ON playlists(owner_id);

CREATE INDEX idx_playlist_items_playlist_position ON playlist_items(playlist_id, position);
CREATE INDEX idx_playlist_items_track_id ON playlist_items(track_id);

CREATE INDEX idx_user_favorites_user_added ON user_favorites(user_id, added DESC);
CREATE INDEX idx_user_favorites_track_id ON user_favorites(track_id);

CREATE INDEX idx_track_history_user_listened ON track_history(user_id, listened_at DESC);
CREATE INDEX idx_track_history_track_id ON track_history(track_id);

CREATE INDEX idx_queue_items_queue_position ON queue_items(queue_id, position);
CREATE INDEX idx_queue_items_track_id ON queue_items(track_id);

CREATE INDEX idx_track_filters_user_id ON track_filters(user_id);
CREATE INDEX idx_api_tokens_user_id ON api_tokens(user_id);
CREATE INDEX idx_user_identities_user_id ON user_identities(user_id);

CREATE INDEX idx_tracks_album_number ON tracks(album_id, number);
