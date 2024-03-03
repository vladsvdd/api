package natsutil

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
)

// NATSClient представляет клиент NATS
type NATSClient struct {
	conn *nats.Conn
}

// NewNATSClient создает новый клиент NATS
func NewNATSClient(url string) (*NATSClient, error) {
	natsConnect, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NATSClient{conn: natsConnect}, nil
}

// Close закрывает соединение с клиентом NATS
func (nc *NATSClient) Close() {
	nc.conn.Close()
}

// Publish отправляет сообщение в NATS
func (nc *NATSClient) Publish(subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return nc.conn.Publish(subject, payload)
}

// MsgHandler адаптирует обычную функцию для обработки сообщений NATS в MsgHandler
type MsgHandler func(msg *nats.Msg)

// SubscribeWithHandler запускает обработчик сообщений NATS с заданным обработчиком.
func (nc *NATSClient) SubscribeWithHandler(subject string, handler MsgHandler) (*nats.Subscription, error) {
	return nc.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg)
	})
}
