Используем команду `goose` для применения этой миграции к нашей базе данных ClickHouse:

```
goose -dir db/migrations/clickhouse/ clickhouse "tcp://test:123456789@localhost:9000/default" status
goose -dir db/migrations/clickhouse/ clickhouse "tcp://test:123456789@localhost:9000/default" up
```

Где:
- `db/migrations/clickhouse/` - это директория, в которой находится файл миграции.
- `tcp://user:password@clickhouse_host:clickhouse_port/db` - это строка подключения к ClickHouse