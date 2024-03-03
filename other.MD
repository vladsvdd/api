Install Nats

# Go client
```go get github.com/nats-io/nats.go/```
# Server
```go get github.com/nats-io/nats-server```

очистка контейнеров
```
docker system prune
```

для запуска локально localhost
1)
```
//natsClient, err := natsutil.NewNATSClient(nats.DefaultURL)
natsClient, err := natsutil.NewNATSClient("nats://nats_hezzl_api:4222")
```
2) в config.yaml
```
host: "postgres_hezzl_api"  # Используется имя контейнера PostgreSQL
#  host: "localhost"
```