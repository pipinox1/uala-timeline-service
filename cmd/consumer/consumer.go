package consumer

import (
	"github.com/nats-io/nats.go"
	"log"
	"uala-timeline-service/config"
)

func SetupConsumer(config *config.Config, deps *config.Dependencies) (*nats.Conn, []*nats.Subscription) {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("Error conectando a NATS:", err)
	}

	qsub, err := nc.QueueSubscribe("post.created", config.ServiceName, handlePostCreated(deps))
	if err != nil {
		log.Fatalf("Error en QueueSubscribe: %v", err)
	}
	qsub2, err := nc.QueueSubscribe("user_timeline.add_post", config.ServiceName, addPostToTimeline(deps))
	if err != nil {
		log.Fatalf("Error en QueueSubscribe: %v", err)
	}
	qsub3, err := nc.QueueSubscribe("user_timeline.remove_post", config.ServiceName, removePostFromTimeline(deps))
	if err != nil {
		log.Fatalf("Error en QueueSubscribe: %v", err)
	}
	subscriptions := []*nats.Subscription{qsub, qsub2, qsub3}
	return nc, subscriptions
}
