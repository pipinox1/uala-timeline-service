package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"time"
	"uala-timeline-service/internal/domain/timeline"
)

var getPostTimelineRow = `
        SELECT 
           post_id,
           published_at
        FROM timelines
        WHERE post_id = $1 AND user_id = $2
    `

var _ timeline.TimelineRepository = (*TimelineRepository)(nil)

type TimelineRepository struct {
	db *sqlx.DB
}

func (t *TimelineRepository) RemovePostFromTimeline(ctx context.Context, userID string, timelinePost timeline.PostTimeline) error {
	//TODO implement me
	panic("implement me")
}

func NewTimelineRepository(db *sqlx.DB) *TimelineRepository {
	return &TimelineRepository{db: db}
}

func (t *TimelineRepository) GetUserPostTimeline(ctx context.Context, userID string, postId string) (*timeline.UserTimeline, error) {
	var postDB postTimelineRow
	err := t.db.GetContext(ctx, &postDB, getPostTimelineRow, postId, userID)
	if err != nil {
		log.Err(err).Msg("error getting user timeline from postgres")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, timeline.ErrUserTimelineNotFound
		}
		return nil, timeline.ErrUserTimelineInternal
	}

	return &timeline.UserTimeline{
		Posts:  []timeline.PostTimeline{postDB.toDomain()},
		UserID: userID,
	}, nil
}

func (t *TimelineRepository) GetUserTimeline(ctx context.Context, userID string, filter timeline.TimelineFilter) (*timeline.UserTimeline, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()

	sb.Select("post_id", "published_at")
	sb.From("timelines")
	sb.Where(sb.Equal("user_id", userID))
	sb.OrderBy("published_at DESC")
	query, args := sb.Build()

	var pgPostTimelineRows []postTimelineRow
	err := t.db.SelectContext(ctx, &pgPostTimelineRows, query, args...)
	if err != nil {
		log.Err(err).Msg("error getting user timeline from postgres")
		return nil, fmt.Errorf("error getting user timeline: %w", err)
	}

	postTimelineRows := make([]timeline.PostTimeline, len(pgPostTimelineRows))
	for i, pgPostTimelineRow := range pgPostTimelineRows {
		postTimelineRows[i] = pgPostTimelineRow.toDomain()
	}

	if len(pgPostTimelineRows) == 0 {
		return &timeline.UserTimeline{
			UserID: userID,
			Posts:  nil,
		}, nil
	}

	return &timeline.UserTimeline{
		UserID: userID,
		Posts:  postTimelineRows,
	}, nil
}

func (t *TimelineRepository) AddPostToUserTimeline(ctx context.Context, userID string, timelinePost timeline.PostTimeline) error {
	_, err := t.db.NamedExecContext(ctx,
		`
		INSERT INTO timelines (user_id, post_id, published_at,created_at)
		VALUES (:user_id, :post_id, :published_at,:created_at)`,
		map[string]interface{}{
			"user_id":      userID,
			"post_id":      timelinePost.PostID,
			"published_at": timelinePost.PublishedAt,
			"created_at":   time.Now(),
		})

	if err != nil {
		log.Err(err).Msg("error adding post to user timeline postgres")
		return err
	}

	return nil
}

type postTimelineRow struct {
	PostID      string    `db:"post_id"`
	PublishedAt time.Time `db:"published_at"`
}

func (p *postTimelineRow) toDomain() timeline.PostTimeline {
	return timeline.PostTimeline{
		PostID:      p.PostID,
		PublishedAt: p.PublishedAt,
	}
}

type TimelineFilter struct {
	Size *int
}
