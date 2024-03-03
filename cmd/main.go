package main

import (
	"api"
	"api/pkg/handler"
	"api/pkg/natsutil"
	"api/pkg/repository"
	"api/pkg/service"
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus" // Пакет logrus для логирования.
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title API goods
// @version 1.0
// @description API goods
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	//Настройка лога файлов
	logFile, err := initLogger()
	if err != nil {
		fmt.Printf("Ошибка при настройке логгера! %v\n", err.Error())
	}
	defer func(logFile *os.File) {
		if err := logFile.Close(); err != nil {
			fmt.Printf("Ошибка при закрытии файла логов! %v\n", err.Error())
		}
	}(logFile)

	// Инициализация конфигурационного файла
	if err := initConfigFile(); err != nil {
		log.Fatal(err)
	}
	log.Println("Конфигурационный файл успешно инициализирован")

	// Загрузка переменных среды из файла .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка при загрузке файла .env: %s", err)
		return
	}
	log.Println("Файл .env успешно загружен")

	// Инициализация подключения к базе данных PostgresSQL
	db, err := initDatabaseConnection()
	if err != nil {
		log.Fatalf("Ошибка инициализации подключения к базе данных: %s", err.Error())
		return
	}
	defer func(db *sqlx.DB) {
		_ = db.Close()
	}(db)
	log.Println("Подключение к базе данных успешно установлено")

	// Инициализация подключения к базе данных PostgresSQL clickhouseDB
	clickhouseDB, err := initClickhouseConnection()
	if err != nil {
		log.Fatalf("Ошибка инициализации подключения к базе данных Clickhouse: %s", err.Error())
		return
	}
	defer func(clickhouseDB *sqlx.DB) {
		_ = clickhouseDB.Close()
	}(clickhouseDB)
	log.Println("Подключение к базе данных Clickhouse успешно установлено")

	redisClient, err := initRedisClient()
	if err != nil {
		log.Fatal("failed to create Redis client: " + err.Error())
		return
	}
	log.Println("Подключение к базе данных Redis успешно установлено")

	// Подключение к серверу Nats
	//natsClient, err := natsutil.NewNATSClient(nats.DefaultURL)
	natsClient, err := natsutil.NewNATSClient("nats://nats_hezzl_api:4222")
	if err != nil {
		log.Fatal("failed to create nats Client : " + err.Error())
		return
	}
	defer natsClient.Close()
	log.Println("Подключение к NATS успешно установлено")

	// Инициализация репозиториев, сервисов и обработчиков
	// ↑3)Работа с БД
	repos := repository.NewRepository(db, clickhouseDB)
	log.Println("Репозитории успешно инициализированы")
	// ↑2)Бизнес логика
	services := service.NewService(repos, redisClient, natsClient)
	log.Println("Бизнес-логика успешно инициализирована")
	// ↑1)Работа с HTTP
	handlers := handler.NewHandler(services)
	log.Println("handlers успешно инициализированы")

	natsGoodsService := service.NewNatsGoodsService(natsClient, repos)
	go natsGoodsService.StartNATSListenerGoodsCreated()
	go natsGoodsService.StartNATSListenerGoodsUpdatedPriority()

	// Ожидание сигналов SIGTERM или SIGINT для завершения работы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	server := new(api.Server)
	// Бесконечный цикл для повторного запуска сервера с интервалом в 5 минут
	for {
		done := make(chan struct{})
		go func() {
			defer close(done)
			if err = server.Run(viper.GetString("port"), handlers.InitRoutes()); !errors.Is(http.ErrServerClosed, err) {
				log.Printf("ОШИБКА при запуске http-сервера: %s\n", err.Error())
			}
		}()
		log.Println("[START] http-сервер успешно запущен")

		select {
		case <-quit:
			// Получен сигнал SIGTERM или SIGINT
			if err = db.Close(); err != nil {
				log.Errorf("ошибка при закрытии соединения с базой данных: %s", err.Error())
			}

			if err = server.Shutdown(context.Background()); err != nil {
				log.Errorf("ошибка при завершении работы сервера: %s", err.Error())
			}
			log.Println("[STOP] http-сервер успешно остановлен")
			return
		case <-done:
			// Вывод сообщения о повторной попытке
			log.Println("[RESTART] Повторная попытка запуска сервера...")

			// Ожидание 5 минут перед следующей попыткой
			time.Sleep(5 * time.Minute)
		}
	}

	select {}
}
