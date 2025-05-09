-- +migrate Up
CREATE TABLE IF NOT EXISTS errors (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL,
    fingerprint CHAR(64) NOT NULL,
    message TEXT NOT NULL,
    stacktrace TEXT NOT NULL,
    file VARCHAR(500) NOT NULL,
    line INTEGER NOT NULL,
    context TEXT,
    time BIGINT NOT NULL,
    created_at INT NOT NULL,
    updated_at INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_errors_fingerprint ON errors(fingerprint);
CREATE INDEX IF NOT EXISTS idx_errors_time ON errors(time);

CREATE TYPE error_groups_status AS ENUM ('unresolved', 'resolved', 'ignored');

CREATE TABLE error_groups (
  fingerprint CHAR(64) PRIMARY KEY,
  project_id UUID NOT NULL,
  type VARCHAR(100) NOT NULL,
  message TEXT NOT NULL,
  file VARCHAR(500) NOT NULL,
  line INTEGER NOT NULL,
  first_seen_at INT NOT NULL,
  last_seen_at INT NOT NULL,
  counter INT DEFAULT 1,
  status error_groups_status DEFAULT 'unresolved'
);

CREATE INDEX IF NOT EXISTS idx_error_groups_project_id ON error_groups(project_id);
CREATE INDEX IF NOT EXISTS idx_error_groups_status ON error_groups(status);