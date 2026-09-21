ALTER TABLE jobs ADD COLUMN unique_key TEXT;

CREATE UNIQUE INDEX idx_jobs_unique_key_active
    ON jobs(unique_key)
    WHERE unique_key IS NOT NULL
      AND unique_key <> ''
      AND status IN ('pending', 'running');