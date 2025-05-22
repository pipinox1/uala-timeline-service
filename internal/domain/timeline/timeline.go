package timeline

import (
	"context"
	"errors"
	"time"
	"uala-timeline-service/internal/domain/posts"
)

var (
	ErrUserTimelineNotFound = errors.New("user_timeline.not_found")
	ErrUserTimelineInternal = errors.New("user_timeline.internal_error")
)

//go:generate mockery --name=TimelineRepository --filename=timeline_follow_repository.go --output=../../../mocks --outpkg=mocks
type TimelineRepository interface {
	GetUserTimeline(ctx context.Context, userID string, filter TimelineFilter) (*UserTimeline, error)
	AddPostToUserTimeline(ctx context.Context, userID string, timelinePost PostTimeline) error
	RemovePostFromTimeline(ctx context.Context, userID string, timelinePost PostTimeline) error
	GetUserPostTimeline(ctx context.Context, userID string, postId string) (*UserTimeline, error)
}

type UserTimeline struct {
	Posts  []PostTimeline
	UserID string
}

type PostTimeline struct {
	PostID      string
	PublishedAt time.Time
}

type TimelineFilter struct {
	DateFrom time.Time
	DateTo   time.Time
	Page     int
}

func CreateTimelinePostFromPost(post posts.Post) PostTimeline {
	return PostTimeline{
		PostID:      post.ID,
		PublishedAt: post.PublishedAt,
	}
}

func CreateUserTimeline(
	userID string,
) (*UserTimeline, error) {
	return &UserTimeline{
		UserID: userID,
		Posts:  []PostTimeline{},
	}, nil
}
