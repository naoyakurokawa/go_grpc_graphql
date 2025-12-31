-- +goose Up
ALTER TABLE tasks
  ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL AFTER updated_at;

-- +goose Down
ALTER TABLE tasks
  DROP COLUMN deleted_at;
