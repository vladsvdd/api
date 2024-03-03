package repository

import (
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
)

type ConfigClickhouse struct {
	Host     string // Хост базы данных
	Port     string // Порт базы данных
	Username string // Имя пользователя базы данных
	Password string // Пароль для подключения к базе данных
	DBName   string // Имя базы данных
}

// NewClickhouseDB создает новое подключение к базе данных Clickhouse на основе переданной конфигурации
func NewClickhouseDB(cfg ConfigClickhouse) (*sqlx.DB, error) {
	db, err := sqlx.Connect("clickhouse", fmt.Sprintf("tcp://%s:%s?username=%s&password=%s&database=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName))
	if err != nil {
		return nil, err
	}

	return db, nil
}
