package application

import (
	"context"
	"fmt"
	"uala-timeline-service/internal/domain/day_timeline_filled"
)

type AddPostToUserTimelineCommand struct {
	UserID string `json:"user_id"`
	PostID string `json:"post_id"`
}

type AddPostToUserTimeline struct {
	timelineService day_timeline_filled.DayUserTimelineFilledService
}

func NewAddPostToUserTimeline(
	timelineService day_timeline_filled.DayUserTimelineFilledService,
) *AddPostToUserTimeline {
	return &AddPostToUserTimeline{
		timelineService: timelineService,
	}
}

func (g *AddPostToUserTimeline) Exec(ctx context.Context, cmd *AddPostToUserTimelineCommand) error {
	fmt.Println("Adding post")
	err := g.timelineService.AddPost(ctx, cmd.PostID, cmd.UserID)
	if err != nil {
		return err
	}

	return nil
}
