ALTER TABLE track_history ADD COLUMN client_event_id TEXT;

-- Backfill existing rows with the web source. New rows must set source
-- explicitly; the API and service reject empty values.
ALTER TABLE track_history ADD COLUMN source TEXT;
UPDATE track_history SET source = 'web';

-- NULL/empty client ids are allowed for old rows and old clients.
CREATE UNIQUE INDEX idx_track_history_user_client_event
    ON track_history(user_id, client_event_id)
    WHERE client_event_id IS NOT NULL AND client_event_id != '';
