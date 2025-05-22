package events

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

type NatsPublisher struct {
	conn *nats.Conn
}

func NewNatsPublisher(host string) *NatsPublisher {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:4222", host))
	if err != nil {
		log.Fatal("Error conectando a NATS:", err)
	}
	return &NatsPublisher{
		conn: nc,
	}
}

func (n *NatsPublisher) Publish(ctx context.Context, event Publishable) error {
	err := n.conn.Publish(event.Topic(), event.Payload())
	if err != nil {
		fmt.Println("Error publishing: " + err.Error())
		return err
	}
	return nil
}

func (n *NatsPublisher) Close() {
	n.conn.Close()
}
