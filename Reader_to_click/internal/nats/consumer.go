package nats

import (
	"fmt"
	nats_pkg "github.com/Georgiy136/go_test/Reader_to_click/pkg/nats"
	"github.com/nats-io/nats.go"
)

type Nats struct {
	nats *nats_pkg.Nats
}

type subjectHandler struct {
	sub *nats.Subscription
}

func NewNats(nats *nats_pkg.Nats) *Nats {
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
