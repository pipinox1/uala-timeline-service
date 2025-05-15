package application

import (
	"context"
	"uala-timeline-service/internal/domain"
)

type RemovePostToUserTimelineTimeCommand struct {
	UserID string
	PostID string
}

type RemovePostToUserTimelineTime struct {
	timelineService domain.TimelineFilledService
}

func NewRemovePostToUserTimelineTime(timelineService domain.TimelineFilledService) *RemovePostToUserTimelineTime {
	return &RemovePostToUserTimelineTime{
		timelineService: timelineService,
	}
}

func (g *RemovePostToUserTimelineTime) Exec(ctx context.Context, cmd *RemovePostToUserTimelineTimeCommand) error {
	err := g.timelineService.AddPost(ctx, cmd.UserID, cmd.PostID)
	if err != nil {
		return err
	}

	return nil
}
