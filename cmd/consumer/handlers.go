package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"uala-timeline-service/config"
	"uala-timeline-service/internal/application"
)

func handlePostCreated(dependencies *config.Dependencies) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		fmt.Println("Recibido")
		log.Printf("Recibido: %s\n", string(msg.Data))
	}
}

func addPostToTimeline(dependencies *config.Dependencies) func(msg *nats.Msg) {
	addPostToTimeline := application.NewAddPostToUserTimeline(dependencies.TimelineService)
	return func(msg *nats.Msg) {
		var cmd application.AddPostToUserTimelineCommand
		err := json.Unmarshal(msg.Data, &cmd)
		if err != nil {
			log.Printf("Error unmarshalling data: %v\n", err)
			return
		}
		addPostToTimeline.Exec(context.Background(), &cmd)
	}
}

func removePostFromTimeline(dependencies *config.Dependencies) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		//TODO implement me
		fmt.Println("removePostFromTimeline")
	}
}
