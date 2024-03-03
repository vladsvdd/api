#!/bin/bash

# Выполнение миграций для PostgreSQL
goose -dir db/migrations/postgres/ postgres "user=postgres dbname=hezzl_api sslmode=disable password=FJIDi3fgnVDsfWE1 host=postgres_hezzl_api" up

# Выполнение миграций для ClickHouse
goose -dir db/migrations/clickhouse/ clickhouse "tcp://test:123456789@clickhouse_hezzl_api:9000/default" up

# Запуск вашего приложения
exec "$@"
