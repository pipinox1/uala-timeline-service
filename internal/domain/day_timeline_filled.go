package domain

import (
	"context"
	"errors"
	"time"
)

type DayUserTimelineFilled struct {
	LastUpdate time.Time
	Posts      []Post
	UserID     string
}

func (t DayUserTimelineFilled) AddPost(post Post) {
	t.Posts = append(t.Posts, post)
}

type DayUserTimelineFilledService interface {
	GetDayUserTimelineFilled(ctx context.Context, filter DayUserTimelineFilledFilter) (*DayUserTimelineFilled, error)
	AddPost(ctx context.Context, postID string, userID string) error
	RemovePost(ctx context.Context, postID string, userID string) error
}

type DayUserTimelineFilledRepository interface {
	GetDayUserTimelineFilled(ctx context.Context, filter DayUserTimelineFilledFilter) (*DayUserTimelineFilled, error)
	AddPosts(ctx context.Context, userID string, post []Post) error
	UpdatePosts(ctx context.Context, userID string, post *Post) error
	RemovePost(ctx context.Context, userID string, post *Post) error
}

type service struct {
	timelineRepository TimelineRepository
	postRepository     PostRepository
	// TODO we use a cache separatly because we will take complex actions in the future to use o refresh the cache
	timelineFilledRepository DayUserTimelineFilledRepository
}

func NewTimelineService(
	timelineRepository TimelineRepository,
	postRepository PostRepository,
	timelineFilledRepository DayUserTimelineFilledRepository,
) DayUserTimelineFilledService {
	return &service{
		timelineRepository:       timelineRepository,
		postRepository:           postRepository,
		timelineFilledRepository: timelineFilledRepository,
	}
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

func (s service) RemovePost(ctx context.Context, postID string, userID string) error {
	post, err := s.postRepository.GetPostById(ctx, postID)
	if err != nil {
		return err
	}

	timelinePost := CreateTimelinePostFromPost(*post)
	err = s.timelineRepository.RemovePostFromTimeline(ctx, userID, timelinePost)
	if err != nil {
		return err
	}

	go func(ctx context.Context) {
		err = s.timelineFilledRepository.RemovePost(ctx, userID, post)
		if err != nil {
			return
		}
	}(context.WithoutCancel(ctx))

	return nil
}

func (s service) GetDayUserTimelineFilled(ctx context.Context, filter DayUserTimelineFilledFilter) (*DayUserTimelineFilled, error) {
	timelineFilled, err := s.timelineFilledRepository.GetDayUserTimelineFilled(ctx, filter)
	if err == nil {
		return timelineFilled, nil
	}

	timeline, err := s.timelineRepository.GetUserTimeline(ctx, filter.UserID, TimelineFilter{
		DateFrom: time.Date(filter.FromYear, time.Month(filter.FromMonth), filter.FromDay, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(filter.ToYear, time.Month(filter.ToMonth), filter.ToDay, 23, 59, 59, 59, time.UTC),
	})
	if err != nil {
		return nil, err
	}

	postIDs := make([]string, len(timeline.Posts))
	for i, post := range timeline.Posts {
		postIDs[i] = post.PostID
	}

	posts, err := s.postRepository.MGetPosts(ctx, postIDs)
	if err != nil {
		return nil, err
	}

	err = s.timelineFilledRepository.AddPosts(ctx, filter.UserID, posts)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s service) AddPost(ctx context.Context, postID string, userID string) error {
	post, err := s.postRepository.GetPostById(ctx, postID)
	if err != nil {
		return err
	}

	_, err = s.timelineRepository.GetUserPostTimeline(ctx, userID, postID)
	if err != nil {
		if errors.Is(err, ErrUserTimelineNotFound) {
			timelinePost := CreateTimelinePostFromPost(*post)
			err = s.timelineRepository.AddPostToUserTimeline(ctx, userID, timelinePost)
			if err != nil {
				return err
			}
		}
	}

	dayTimeline, err := s.timelineFilledRepository.GetDayUserTimelineFilled(ctx, DayUserTimelineFilledFilter{
		UserID:    userID,
		FromDay:   post.PublishedAt.Day(),
		FromMonth: int(post.PublishedAt.Month()),
		FromYear:  post.PublishedAt.Year(),
	})
	if err != nil {
		return err
	}

	for _, dayTimelinePost := range dayTimeline.Posts {
		if dayTimelinePost.ID == post.ID {
			//Discard add legacy message
			if post.UpdatedAt.Before(dayTimelinePost.UpdatedAt) || post.UpdatedAt.Equal(dayTimelinePost.UpdatedAt) {
				return nil
			}
			err = s.timelineFilledRepository.UpdatePosts(ctx, userID, post)
			if err != nil {
				return err
			}
			return nil
		}
	}

	err = s.timelineFilledRepository.AddPosts(ctx, userID, []Post{*post})
	if err != nil {
		return err
	}

	return nil
}
