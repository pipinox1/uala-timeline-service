package day_timeline_filled

import (
	"context"
	"time"
	"uala-timeline-service/internal/domain/posts"
)

//go:generate mockery --name=DayUserTimelineFilledRepository --filename=mocks_day_timeline_filled_repository.go --output=../../../mocks --outpkg=mocks
type DayUserTimelineFilledRepository interface {
	GetDayUserTimelineFilled(ctx context.Context, filter DayUserTimelineFilledFilter) (*DayUserTimelineFilled, error)
	AddPosts(ctx context.Context, userID string, post []posts.Post) error
	UpdatePosts(ctx context.Context, userID string, post *posts.Post) error
	RemovePost(ctx context.Context, userID string, post *posts.Post) error
}

type DayUserTimelineFilled struct {
	LastUpdate time.Time
	Posts      []posts.Post
	UserID     string
}

type DayUserTimelineFilledFilter struct {
	UserID    string
	FromDay   int
	FromMonth int
	FromYear  int
	ToMonth   int
	ToYear    int
	ToDay     int
	Page      int
}

func (t DayUserTimelineFilled) AddPost(post posts.Post) {
	t.Posts = append(t.Posts, post)
}

func CreateDayUserTimelineFilled(userID string, posts []posts.Post) DayUserTimelineFilled {
	return DayUserTimelineFilled{
		LastUpdate: time.Now(),
		Posts:      posts,
		UserID:     userID,
	}
}
