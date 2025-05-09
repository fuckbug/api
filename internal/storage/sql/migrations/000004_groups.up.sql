-- +migrate Up

alter table error_groups rename column fingerprint to id;
alter table error_groups drop column type;

CREATE TYPE log_groups_status AS ENUM ('unresolved', 'resolved', 'ignored');

CREATE TABLE log_groups (
  id CHAR(64) PRIMARY KEY,
  project_id UUID NOT NULL,
  level log_level NOT NULL,
  message TEXT NOT NULL,
  first_seen_at INT NOT NULL,
  last_seen_at INT NOT NULL,
  counter INT DEFAULT 1,
  status error_groups_status DEFAULT 'unresolved'
);

CREATE INDEX IF NOT EXISTS idx_log_groups_project_id ON log_groups(project_id);
CREATE INDEX IF NOT EXISTS idx_log_groups_status ON log_groups(status);