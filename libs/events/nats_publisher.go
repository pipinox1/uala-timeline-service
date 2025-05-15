package events

import (
	"context"
	"github.com/nats-io/nats.go"
	"log"
)

type NatsPublisher struct {
	nc *nats.Conn
}

func NewNatsPublisher() *NatsPublisher {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("Error conectando a NATS:", err)
	}
	defer nc.Close()

	return &NatsPublisher{}
}

func (n *NatsPublisher) Publish(ctx context.Context, event Publishable) error {
	return n.nc.Publish(event.Topic(), event.Payload())
}
