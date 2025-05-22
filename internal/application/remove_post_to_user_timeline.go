package application

import (
	"context"
	"uala-timeline-service/internal/domain/day_timeline_filled/service"
)

type RemovePostToUserTimelineTimeCommand struct {
	UserID string
	PostID string
}

type RemovePostToUserTimelineTime struct {
	timelineService service.DayUserTimelineFilledService
}

func NewRemovePostToUserTimelineTime(timelineService service.DayUserTimelineFilledService) *RemovePostToUserTimelineTime {
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
