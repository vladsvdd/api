package api

import (
	"context"  // Импорт пакета context для работы с контекстами в Go.
	"net/http" // Импорт пакета http для создания HTTP-сервера.
	"time"     // Импорт пакета time для работы со временем.
)

// Server Определение структуры, которая содержит http-сервер
type Server struct {
	httpServer *http.Server
}

// Run Метод используется для запуска http-сервера на указанном порту с определенным обработчиком
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{ // Инициализация http-сервера
		Addr:           ":" + port,       // Указание порта для прослушивания
		Handler:        handler,          // Установка обработчика запросов
		MaxHeaderBytes: 1 << 20,          // Установка максимального размера заголовка запроса (1 MB)
		ReadTimeout:    10 * time.Second, // Установка времени ожидания на чтение
		WriteTimeout:   10 * time.Second, // Установка времени ожидания на запись
	}

	return s.httpServer.ListenAndServe() // Запуск http-сервера для прослушивания входящих запросов.
}

// Shutdown Метод используется для корректного завершения работы http-сервера с использованием контекста
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx) // Завершение работы http-сервера с использованием контекста.
}
