package main

import (
	"api/pkg/repository"
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

const (
	LogFileName = "./logger/log.logs"
	ConfigName  = "config"
	ConfigPath  = "configs"
	ConfigType  = "yaml"
)

// initConfigFile используется для инициализации конфигурационного файла
func initConfigFile() error {
	viper.SetConfigName(ConfigName)
	viper.AddConfigPath(ConfigPath)
	viper.SetConfigType(ConfigType) // Добавьте эту строку, если файл имеет расширение .yaml
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("ошибка инициализации конфиг файла %s", err.Error())
	}

	return nil
}

// initDatabaseConnection используется для инициализации конфигурационного файла
func initDatabaseConnection() (*sqlx.DB, error) {
	return repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	})
}

// initClickhouseConnection используется для инициализации конфигурационного файла
func initClickhouseConnection() (*sqlx.DB, error) {
	return repository.NewClickhouseDB(repository.ConfigClickhouse{
		Host:     viper.GetString("clickhouse.host"),
		Port:     viper.GetString("clickhouse.port"),
		Username: viper.GetString("clickhouse.username"),
		DBName:   viper.GetString("clickhouse.dbname"),
		Password: os.Getenv("CLICKHOUSE_PASSWORD"),
	})
}

// initDatabaseConnection используется для инициализации конфигурационного файла
func initRedisClient() (*repository.RedisClient, error) {
	redisAddress := fmt.Sprintf("%s:%s",
		viper.GetString("redis.host"),
		viper.GetString("redis.port")) //"localhost:6379"
	redisPassword := "" //os.Getenv("REDIS_PASSWORD")
	return repository.NewRedisClient(redisAddress, redisPassword, 0)
}

// initLogger создает логгер и настраивает форматирование журнала с использованием log
func initLogger() (*os.File, error) {
	file, err := os.OpenFile(LogFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл журнала: %s", err)
	}

	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{})

	return file, nil
}
