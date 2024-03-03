-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS goods
(
    id          SERIAL PRIMARY KEY,
    project_id  BIGINT    NOT NULL,
    name        TEXT      NOT NULL,
    description TEXT,
    priority    INT       NOT NULL,
    removed     BOOL               DEFAULT FALSE NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_goods_project_id FOREIGN KEY (project_id) REFERENCES projects (id)
);

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_goods_project_id ON goods (project_id);
CREATE INDEX IF NOT EXISTS idx_goods_name ON goods (name);

-- Add trigger function
CREATE OR REPLACE FUNCTION set_priority()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.priority := (SELECT COALESCE(MAX(priority), 0) + 1 FROM goods);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Add trigger
CREATE TRIGGER set_priority_trigger
    BEFORE INSERT
    ON goods
    FOR EACH ROW
EXECUTE FUNCTION set_priority();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS goods;
-- +goose StatementEnd


