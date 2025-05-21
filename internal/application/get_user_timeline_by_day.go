package application

import (
	"context"
	"errors"
	"uala-timeline-service/internal/domain"
)

var (
	DatesFieldAreMandatory = errors.New("dates fields are mandatory")
)

type GetUserTimelineCommand struct {
	UserID    string `json:"-"`
	FromDay   int    `json:"from_day"`
	FromMonth int    `json:"from_month"`
	FromYear  int    `json:"from_year"`
	ToDay     int    `json:"to_day"`
	ToMonth   int    `json:"to_month"`
	ToYear    int    `json:"to_year"`
}

type GetUserTimelineResponse struct {
	*TimelineFilled
}

type GetUserTimeline struct {
	timelineService domain.DayUserTimelineFilledService
}

func NewGetUserTimeline(timelineService domain.DayUserTimelineFilledService) *GetUserTimeline {
	return &GetUserTimeline{
		timelineService: timelineService,
	}
}

func (g *GetUserTimeline) Exec(ctx context.Context, cmd *GetUserTimelineCommand) (*GetUserTimelineResponse, error) {
	if cmd.ToDay == 0 || cmd.ToMonth == 0 || cmd.ToYear == 0 || cmd.FromDay == 0 || cmd.FromMonth == 0 || cmd.FromYear == 0 {
		return nil, DatesFieldAreMandatory
	}
	userTimeline, err := g.timelineService.GetDayUserTimelineFilled(ctx, domain.DayUserTimelineFilledFilter{
		UserID:    cmd.UserID,
		FromDay:   cmd.FromDay,
		FromMonth: cmd.FromMonth,
		FromYear:  cmd.FromYear,
		ToMonth:   cmd.ToDay,
		ToYear:    cmd.ToMonth,
		ToDay:     cmd.ToYear,
	})
	if err != nil {
		return nil, err
	}

	return &GetUserTimelineResponse{
		FromDomain(userTimeline),
	}, nil
}
