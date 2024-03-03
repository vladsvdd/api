package service

import (
	"api/models"
	"api/pkg/natsutil"
	"api/pkg/repository"
	"encoding/json"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

type GoodsNatsService struct {
	repo       *repository.Repository
	natsClient *natsutil.NATSClient
}

func NewNatsGoodsService(
	natsClient *natsutil.NATSClient,
	repo *repository.Repository,
) *GoodsNatsService {
	return &GoodsNatsService{
		natsClient: natsClient,
		repo:       repo,
	}
}

// StartNATSListenerGoodsCreated запускает слушателя сообщений NATS
func (s *GoodsNatsService) StartNATSListenerGoodsCreated() {
	messageChan := make(chan models.Goods, 4) // Буферизованный канал для хранения записей

	// Запускаем горутину для обработки записей из канала
	go func() {
		var buffer []models.Goods
		for {
			select {
			case good := <-messageChan:
				buffer = append(buffer, good)
				if len(buffer) >= 4 { // Проверяем, достигли ли мы нужного количества записей
					if err := s.repo.CreateBatch(buffer); err != nil {
						log.Error(err)
						// Возможно, вам нужно обработать ошибку, например, повторно отправить записи в канал
					}
					buffer = nil
				}
			}
		}
	}()

	_, err := s.natsClient.SubscribeWithHandler(GoodsNatsKey, s.handleNATSMessage(messageChan))
	if err != nil {
		log.Fatal(err)
	}
}

// StartNATSListenerGoodsUpdatedPriority запускает слушателя сообщений NATS
func (s *GoodsNatsService) StartNATSListenerGoodsUpdatedPriority() {
	messageChan := make(chan []models.GoodsPriority)

	_, err := s.natsClient.SubscribeWithHandler(GoodsUpdatedPriorityNatsKey, func(msg *nats.Msg) {
		log.Printf("Received NATS message: %s\n", string(msg.Data))

		var goods []models.GoodsPriority
		if err := json.Unmarshal(msg.Data, &goods); err != nil {
			log.Error(err)
			return
		}

		messageChan <- goods
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		for {
			select {
			case goods := <-messageChan:
				if err := s.repo.CreatePriorityBatch(goods); err != nil {
					log.Error(err)
					// TODO: можно повторно отправить записи в канал
				}
			}
		}
	}()
}

// handleNATSMessage обрабатывает сообщения NATS
func (s *GoodsNatsService) handleNATSMessage(messageChan chan models.Goods) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		log.Printf("Received NATS message: %s\n", string(msg.Data))

		var good models.Goods
		if err := json.Unmarshal(msg.Data, &good); err != nil {
			log.Error(err)
			return
		}

		messageChan <- good // Отправляем запись в буферизованный канал
	}
}
