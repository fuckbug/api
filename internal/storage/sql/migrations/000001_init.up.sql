-- +migrate Up
CREATE TYPE log_level AS ENUM ('INFO', 'WARN', 'ERROR', 'DEBUG');

CREATE TABLE IF NOT EXISTS logs (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL,
    level log_level NOT NULL,
    message TEXT NOT NULL,
    context TEXT,
    time BIGINT NOT NULL,
    created_at INT NOT NULL,
    updated_at INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_logs_time ON logs(time);
CREATE INDEX IF NOT EXISTS idx_logs_level ON logs(level);