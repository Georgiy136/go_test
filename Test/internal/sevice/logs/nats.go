package logs

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"myapp/config"
	"myapp/internal/models"
	"myapp/pkg/nats"
)

type NatsLogging struct {
	nats nats.Nats
	cfg  config.Nats
}

func NewNatsLogging(nats nats.Nats, cfg config.Nats) *NatsLogging {
	return &NatsLogging{
		nats: nats,
		cfg:  cfg,
	}
}

func (n *NatsLogging) SendLogToNats(data models.Log) error {
	bytesData, err := jsoniter.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshal log data: %w, data: %v", err, data)
	}

	if _, err = n.nats.Js.Publish(n.cfg.ChannelName, bytesData); err != nil {
		return fmt.Errorf("error publish to stream: %w, name: %s", err, n.cfg.ChannelName)
	}
	return nil
}
