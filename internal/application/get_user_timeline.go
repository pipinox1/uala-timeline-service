package application

import (
	"context"
	"errors"
	"uala-timeline-service/internal/domain"
)

type GetUserTimelineCommand struct {
	UserID string
}

type GetUserTimelineResponse struct {
	*TimelineFilled
}

type GetUserTimeline struct {
	timelineService domain.TimelineFilledService
}

func NewGetUserTimeline(timelineService domain.TimelineFilledService) *GetUserTimeline {
	return &GetUserTimeline{
		timelineService: timelineService,
	}
}

func (g *GetUserTimeline) Exec(ctx context.Context, cmd *GetUserTimelineCommand) (*GetUserTimelineResponse, error) {
	userTimeline, err := g.timelineService.GetUserTimeline(ctx, cmd.UserID, domain.TimelineFilter{})
	if err != nil {
		if !errors.Is(err, domain.ErrTimelineNotFound) {
			return nil, err
		}
		userTimeline, err = g.timelineService.CreateUserTimeline(ctx, cmd.UserID)
		if err != nil {
			return nil, err
		}
	}

	return &GetUserTimelineResponse{
		FromDomain(userTimeline),
	}, nil
}
