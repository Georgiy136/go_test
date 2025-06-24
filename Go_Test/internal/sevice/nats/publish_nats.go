package nats

import (
	"fmt"
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
	if n.nats != nil {
		if _, err := n.nats.Js.Publish(channelName, data); err != nil {
			return fmt.Errorf("error publish to stream: %w, name: %s", err, channelName)
		}
	}
	return nil
}
