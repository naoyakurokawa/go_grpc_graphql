-- +goose Up
CREATE TABLE tasks (
   id varchar(255) NOT NULL,
   title varchar(255) DEFAULT NULL,
   note text DEFAULT NULL,
   completed integer DEFAULT 0,
   created_at TIMESTAMP DEFAULT NULL,
   updated_at TIMESTAMP DEFAULT NULL
);

-- +goose Down
DROP TABLE tasks;