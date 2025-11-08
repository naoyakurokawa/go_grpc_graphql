-- +goose Up
ALTER TABLE tasks
  ADD COLUMN due_date DATE NULL AFTER completed;

-- +goose Down
ALTER TABLE tasks
  DROP COLUMN due_date;
