-- +migrate Up
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    public_key TEXT NOT NULL,
    created_at INT NOT NULL,
    updated_at INT NOT NULL,
    deleted_at INT DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
CREATE INDEX IF NOT EXISTS idx_projects_deleted_at ON projects(deleted_at);