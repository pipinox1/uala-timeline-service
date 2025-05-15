package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTimelineFilledNotFound = errors.New("timeline.not_found")
)

type TimelineFilled struct {
	LastUpdate time.Time
	Posts      []Post
	UserID     string
}

func (t TimelineFilled) AddPost(post Post) {
	t.Posts = append(t.Posts, post)
}

type TimelineFilledService interface {
	CreateUserTimeline(ctx context.Context, userID string) (*TimelineFilled, error)
	GetUserTimeline(ctx context.Context, userID string, filter TimelineFilter) (*TimelineFilled, error)
	AddPost(ctx context.Context, postID string, userID string) error
	RemovePost(ctx context.Context, postID string, userID string) error
}

type TimelineFilledRepository interface {
	GetUserTimeline(ctx context.Context, userID string) (*TimelineFilled, error)
	Update(ctx context.Context, timelineFilled *TimelineFilled) error
}

type service struct {
	timelineRepository TimelineRepository
	postRepository     PostRepository
	// TODO we use a cache separatly because we will take complex actions in the future to use o refresh the cache
	timelineFilledRepository TimelineFilledRepository
}

func NewTimelineService(
	timelineRepository TimelineRepository,
	postRepository PostRepository,
	timelineFilledRepository TimelineFilledRepository,
) TimelineFilledService {
	return &service{
		timelineRepository:       timelineRepository,
		postRepository:           postRepository,
		timelineFilledRepository: timelineFilledRepository,
	}
}

// TODO improve this method to receive a post because it ocurrs in the add post to timeline
func (s service) CreateUserTimeline(ctx context.Context, userID string) (*TimelineFilled, error) {
	createdTimeline, err := CreateUserTimeline(userID)
	if err != nil {
		return nil, err
	}
	err = s.timelineRepository.SaveTimeline(ctx, createdTimeline)
	if err != nil {
		return nil, err
	}
	return &TimelineFilled{
		LastUpdate: createdTimeline.LastUpdate,
		Posts:      nil,
		UserID:     createdTimeline.UserID,
	}, nil
}

func (s service) RemovePost(ctx context.Context, postID string, userID string) error {
	//TODO implement me
	panic("implement me")
}

func (s service) GetUserTimeline(ctx context.Context, userID string, filter TimelineFilter) (*TimelineFilled, error) {
	timelineFilled, err := s.timelineFilledRepository.GetUserTimeline(ctx, userID)
	if err == nil {
		return timelineFilled, nil
	}

	timeline, err := s.timelineRepository.GetUserTimeline(ctx, userID, TimelineFilter{})
	if err != nil {
		return nil, err
	}

	posts, err := s.postRepository.MGetPosts(ctx, timeline.Posts)
	if err != nil {
		return nil, err
	}

	timelineFilled = &TimelineFilled{
		LastUpdate: timeline.LastUpdate,
		Posts:      posts,
		UserID:     timeline.UserID,
	}

	err = s.timelineFilledRepository.Update(ctx, timelineFilled)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s service) AddPost(ctx context.Context, postID string, userID string) error {
	err := s.timelineRepository.AddPostToTimeline(ctx, userID, postID)
	if err != nil {
		return err
	}

	go func(ctx context.Context) {
		timelineFilled, err := s.timelineFilledRepository.GetUserTimeline(ctx, userID)
		if err != nil {
			return
		}

		post, err := s.postRepository.GetPostById(ctx, postID)
		if err != nil {
			return
		}

		timelineFilled.AddPost(*post)
		err = s.timelineFilledRepository.Update(ctx, timelineFilled)
		if err != nil {
			return
		}
	}(context.WithoutCancel(ctx))

	return nil
}
