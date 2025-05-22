package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"uala-timeline-service/config"
	"uala-timeline-service/internal/application"
)

func handlePostCreated(dependencies *config.Dependencies) func(msg *nats.Msg) {
	splitPostUpdateForUsers := application.NewSplitPostUpdateForUsers(
		dependencies.PostRepository,
		dependencies.FollowRepository,
		dependencies.EventPublisher,
	)
	return func(msg *nats.Msg) {
		log.Info().Msg("handlePostCreated event")
		var cmd application.SplitPostUpdateForUsersCommand
		err := json.Unmarshal(msg.Data, &cmd)
		if err != nil {
			log.Err(err)
			return
		}
		err = splitPostUpdateForUsers.Exec(context.Background(), &cmd)
		if err != nil {
			msg.Nak()
		}
		msg.Ack()
	}
}

func addPostToTimeline(dependencies *config.Dependencies) func(msg *nats.Msg) {
	addPostToTimeline := application.NewAddPostToUserTimeline(dependencies.TimelineService)
	return func(msg *nats.Msg) {
		log.Info().Msg("addPostToTimeline event")
		var cmd application.AddPostToUserTimelineCommand
		err := json.Unmarshal(msg.Data, &cmd)
		if err != nil {
			log.Err(err)
			return
		}
		err = addPostToTimeline.Exec(context.Background(), &cmd)
		if err != nil {
			msg.Nak()
		}
		msg.Ack()
	}
}

func removePostFromTimeline(dependencies *config.Dependencies) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		//TODO implement me
		fmt.Println("removePostFromTimeline")
	}
}
