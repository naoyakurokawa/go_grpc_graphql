-- +goose Up
CREATE TABLE categories (
   id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
   name VARCHAR(255) NOT NULL,
   created_at TIMESTAMP NULL DEFAULT NULL,
   updated_at TIMESTAMP NULL DEFAULT NULL
);

ALTER TABLE tasks
ADD COLUMN category_id BIGINT UNSIGNED NULL AFTER note,
ADD CONSTRAINT fk_tasks_category_id FOREIGN KEY (category_id) REFERENCES categories(id);

INSERT INTO categories (name) VALUES ('Work'), ('Personal'), ('Other');

-- +goose Down
ALTER TABLE tasks
DROP FOREIGN KEY fk_tasks_category_id,
DROP COLUMN category_id;

DROP TABLE categories;
