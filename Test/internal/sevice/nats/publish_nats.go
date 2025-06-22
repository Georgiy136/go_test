package nats

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"myapp/pkg/nats"
)

type NatsService struct {
	nats *nats.Nats
}

func NewNatsService(nats *nats.Nats) *NatsService {
	return &NatsService{
		nats: nats,
	}
}

func (n *NatsService) SendBatch(channelName string, data []byte) error {
	bytesData, err := jsoniter.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshal log data: %w, data: %v", err, data)
	}

	if _, err = n.nats.Js.Publish(channelName, bytesData); err != nil {
		return fmt.Errorf("error publish to stream: %w, name: %s", err, channelName)
	}
	return nil
}
