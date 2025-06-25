package nats

import (
	"fmt"
	"github.com/Georgiy136/go_test/Reader_to_click/pkg/nats"
)

type Nats struct {
	nats *nats.Nats
}

func NewNats(nats *nats.Nats) *Nats {
	return &Nats{
		nats: nats,
	}
}

func (n *Nats) GetData(channelName string) ([]byte, error) {
	msg, err := n.nats.Js.GetMsg(channelName, 1)
	if err != nil {
		return nil, fmt.Errorf("error GetData, stream: %w, err: %s", channelName, err)
	}
	return msg.Data, nil
}
