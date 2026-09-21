ALTER TABLE jobs ADD COLUMN next_attempt_at INTEGER NOT NULL DEFAULT 0;

CREATE INDEX idx_jobs_status_next_attempt ON jobs(status, next_attempt_at);