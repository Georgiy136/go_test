package nats

import (
	"fmt"
	"github.com/Georgiy136/go_test/go_test/pkg/nats"
)

type NatsService struct {
	nats *nats.Nats
}

func NewNatsService(nats *nats.Nats) *NatsService {
	return &NatsService{
		nats: nats,
	}
}

func (n *NatsService) SendBatch(data []byte) error {
	if n.nats != nil {
		if _, err := n.nats.Js.Publish(n.nats.Cfg.ChannelName, data); err != nil {
			return fmt.Errorf("error publish to stream: %w, name: %s", err, n.nats.Cfg.ChannelName)
		}
	}
	return nil
}
