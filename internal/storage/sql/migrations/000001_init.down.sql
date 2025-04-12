-- +migrate Down
DROP TABLE IF EXISTS logs;
DROP TYPE IF EXISTS log_level;