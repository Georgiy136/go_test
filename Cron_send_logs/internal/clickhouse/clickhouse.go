package nats

import (
	"fmt"
	"github.com/Georgiy136/go_test/Cron_send_logs/pkg/nats"
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
	subj, err := n.nats.Js.GetMsg(channelName, 1)
	if err != nil {
		return nil, fmt.Errorf("error GetData from stream: %s, err: %s", channelName, err)
	}

	return subj.Data, nil
}
