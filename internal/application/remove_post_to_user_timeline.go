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
	timelineService domain.DayUserTimelineFilledService
}

func NewRemovePostToUserTimelineTime(timelineService domain.DayUserTimelineFilledService) *RemovePostToUserTimelineTime {
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
