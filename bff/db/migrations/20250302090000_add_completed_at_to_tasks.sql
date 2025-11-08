-- +goose Up
ALTER TABLE tasks
  ADD COLUMN completed_at DATETIME NULL AFTER completed;

-- +goose Down
ALTER TABLE tasks
  DROP COLUMN completed_at;
