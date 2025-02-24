-- +goose Up
CREATE TABLE tasks (
   id INTEGER AUTO_INCREMENT PRIMARY KEY,
   title varchar(255) DEFAULT NULL,
   note text DEFAULT NULL,
   completed integer DEFAULT 0,
   created_at TIMESTAMP DEFAULT NULL,
   updated_at TIMESTAMP DEFAULT NULL
);

-- +goose Down
DROP TABLE tasks;