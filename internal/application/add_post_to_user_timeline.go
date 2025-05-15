package application

import (
	"context"
	"uala-timeline-service/internal/domain"
)

type AddPostToUserTimelineCommand struct {
	UserID string `json:"user_id"`
	PostID string `json:"post_id"`
}

type AddPostToUserTimeline struct {
	timelineService domain.TimelineFilledService
}

func NewAddPostToUserTimeline(
	timelineService domain.TimelineFilledService,
) *AddPostToUserTimeline {
	return &AddPostToUserTimeline{
		timelineService: timelineService,
	}
}

func (g *AddPostToUserTimeline) Exec(ctx context.Context, cmd *AddPostToUserTimelineCommand) error {
	err := g.timelineService.AddPost(ctx, cmd.UserID, cmd.PostID)
	if err != nil {
		return err
	}

	return nil
}
