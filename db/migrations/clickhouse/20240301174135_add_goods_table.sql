-- +goose Up
-- +goose StatementBegin
-- Создаем таблицу goods
CREATE TABLE IF NOT EXISTS goods
(
    id          Int64,
    ProjectId   Int64,
    Name        String,
    Description String,
    Priority    Int32,
    Removed     UInt8 DEFAULT 0,
    EventTime   DateTime DEFAULT now()
) ENGINE = MergeTree()
      ORDER BY id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Удаляем таблицу goods
DROP TABLE IF EXISTS goods;
-- +goose StatementEnd
