package domain

import (
	"context"
	"errors"
	"time"
)

var ErrTimelineNotFound = errors.New("timeline.not_found")

type TimelineRepository interface {
	GetUserTimeline(ctx context.Context, userID string, filter TimelineFilter) (*UserTimeline, error)
	SaveTimeline(ctx context.Context, timeline *UserTimeline) error
	AddPostToTimeline(ctx context.Context, userID string, postID string) error
}

type UserTimeline struct {
	LastUpdate time.Time
	Posts      []string
	UserID     string
}

type TimelineFilter struct {
	Size *int
}

func CreateUserTimeline(
	userID string,
) (*UserTimeline, error) {
	return &UserTimeline{
		UserID:     userID,
		Posts:      []string{},
		LastUpdate: time.Now(),
	}, nil
}
