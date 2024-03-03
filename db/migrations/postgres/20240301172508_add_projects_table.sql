-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS projects
(
    id         SERIAL PRIMARY KEY,
    name       TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Insert initial data
INSERT INTO projects (name)
VALUES ('Первая запись');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS projects;
-- +goose StatementEnd