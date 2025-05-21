package infrastructure

/*
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
	"uala-timeline-service/internal/domain"
)

var _ domain.DayUserTimelineFilledRepository = (*TimelineFilledRepository)(nil)

type TimelineFilledRepository struct {
	redisClient *redis.Client
	expiration  time.Duration
}

func NewTimelineFilledRepository(redisClient *redis.Client) *TimelineFilledRepository {
	return &TimelineFilledRepository{
		redisClient: redisClient,
		expiration:  1 * time.Hour,
	}
}

func (t *TimelineFilledRepository) GetUserTimeline(ctx context.Context, userID string) (*domain.DayUserTimelineFilled, error) {
	key := buildTimelineKey(userID)

	data, err := t.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domain.ErrTimelineFilledNotFound
		}
		return nil, fmt.Errorf("error getting timeline from redis: %w", err)
	}

	var redisTimeline redisTimelineFilled
	if err := json.Unmarshal(data, &redisTimeline); err != nil {
		return nil, fmt.Errorf("error unmarshaling timeline data: %w", err)
	}

	timeline := mapRedisToDomain(redisTimeline)

	return timeline, nil
}

func (t *TimelineFilledRepository) Update(ctx context.Context, timelineFilled *domain.DayUserTimelineFilled) error {
	key := buildTimelineKey(timelineFilled.UserID)
	data, err := json.Marshal(mapDomainToRedis(timelineFilled))
	if err != nil {
		return fmt.Errorf("error marshaling timeline data: %w", err)
	}

	if err := t.redisClient.Set(ctx, key, data, t.expiration).Err(); err != nil {
		return fmt.Errorf("error setting timeline in redis: %w", err)
	}

	return nil
}

func buildTimelineKey(userID string) string {
	return fmt.Sprintf("timeline:%s", userID)
}

func mapDomainToRedis(timeline *domain.DayUserTimelineFilled) redisTimelineFilled {
	redisTimeline := redisTimelineFilled{
		LastUpdate: timeline.LastUpdate,
		UserID:     timeline.UserID,
	}

	posts := make([]redisPost, len(timeline.Posts))
	for i, post := range timeline.Posts {
		contents := make([]redisContent, len(post.Contents))
		for i, content := range post.Contents {
			contents[i] = redisContent{
				Type: content.Type,
				Text: content.Text,
				Url:  content.Url,
			}
		}

		rPost := redisPost{
			ID:       post.ID,
			AuthorID: post.AuthorID,
			Contents: make([]redisContent, 0, len(post.Contents)),
		}
		posts[i] = rPost
	}

	return redisTimeline
}

func mapRedisToDomain(redisTimeline redisTimelineFilled) *domain.DayUserTimelineFilled {
	timeline := &domain.DayUserTimelineFilled{
		LastUpdate: redisTimeline.LastUpdate,
		UserID:     redisTimeline.UserID,
		Posts:      make([]domain.Post, 0, len(redisTimeline.Posts)),
	}

	for _, rPost := range redisTimeline.Posts {
		post := domain.Post{
			ID:       rPost.ID,
			AuthorID: rPost.AuthorID,
			Contents: make([]domain.Content, 0, len(rPost.Contents)),
		}

		for _, rContent := range rPost.Contents {
			content := domain.Content{
				Type: rContent.Type,
				Text: rContent.Text,
				Url:  rContent.Url,
			}
			post.Contents = append(post.Contents, content)
		}

		timeline.Posts = append(timeline.Posts, post)
	}

	return timeline
}

type redisTimelineFilled struct {
	LastUpdate time.Time   `json:"last_update"`
	Posts      []redisPost `json:"posts"`
	UserID     string      `json:"user_id"`
}

type redisPost struct {
	ID       string         `json:"id"`
	Contents []redisContent `json:"contents"`
	AuthorID string         `json:"author_id"`
}

type redisContent struct {
	Type string  `json:"type"`
	Text *string `json:"text,omitempty"`
	Url  *string `json:"url,omitempty"`
}

*/
