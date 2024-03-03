# Используйте официальный образ Golang
FROM golang:1.21.3

# Установка переменных среды из .env файла
ARG GO_API_PORT_REMOTE
ENV GO_API_PORT_REMOTE=${GO_API_PORT_REMOTE}

# Создайте директорию приложения внутри контейнера
WORKDIR /app

# Скопируйте файлы go.mod и go.sum внутрь контейнера
COPY go.mod .
COPY go.sum .

# Загрузите зависимости с помощью команды go mod download
RUN go mod download

# Удалите все файлы из директории ./go в контейнере
RUN rm -rf ./go/*

RUN echo "module ${DOMAIN}/api" > go.mod

# Скопируйте только нужные файлы из . в директорию ./go в контейнере
COPY . .

# Устанавливаем goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Сборка приложения
RUN go build -o api ./cmd

# Откройте порт, на котором будет работать приложение
EXPOSE ${GO_API_PORT_REMOTE}

# Запуск приложения
CMD ["./api"]