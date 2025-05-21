package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"time"
	"uala-timeline-service/internal/domain"
)

var getPostTimelineRow = `
        SELECT 
           post_id,
           published_at
        FROM timelines
        WHERE post_id = $1 AND user_id = $2
    `

var _ domain.TimelineRepository = (*TimelineRepository)(nil)

type TimelineRepository struct {
	db *sqlx.DB
}

func (t *TimelineRepository) RemovePostFromTimeline(ctx context.Context, userID string, timelinePost domain.PostTimeline) error {
	//TODO implement me
	panic("implement me")
}

func NewTimelineRepository(db *sqlx.DB) *TimelineRepository {
	return &TimelineRepository{db: db}
}

func (t *TimelineRepository) GetUserPostTimeline(ctx context.Context, userID string, postId string) (*domain.UserTimeline, error) {
	var postDB postTimelineRow
	err := t.db.GetContext(ctx, &postDB, getPostTimelineRow, postId, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserTimelineNotFound
		}
		return nil, domain.ErrUserTimelineInternal
	}

	return &domain.UserTimeline{
		Posts:  []domain.PostTimeline{postDB.toDomain()},
		UserID: userID,
	}, nil
}

func (t *TimelineRepository) GetUserTimeline(ctx context.Context, userID string, filter domain.TimelineFilter) (*domain.UserTimeline, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()

	sb.Select("post_id", "published_at")
	sb.From("timelines")
	sb.Where(sb.Equal("user_id", userID))
	sb.OrderBy("published_at DESC")
	query, args := sb.Build()

	var pgPostTimelineRows []postTimelineRow
	err := t.db.SelectContext(ctx, &pgPostTimelineRows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting user timeline: %w", err)
	}

	postTimelineRows := make([]domain.PostTimeline, len(pgPostTimelineRows))
	for i, pgPostTimelineRow := range pgPostTimelineRows {
		postTimelineRows[i] = pgPostTimelineRow.toDomain()
	}

	if len(pgPostTimelineRows) == 0 {
		return &domain.UserTimeline{
			UserID: userID,
			Posts:  nil,
		}, nil
	}

	return &domain.UserTimeline{
		UserID: userID,
		Posts:  postTimelineRows,
	}, nil
}

func (t *TimelineRepository) AddPostToUserTimeline(ctx context.Context, userID string, timelinePost domain.PostTimeline) error {
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
		return err
	}

	return nil
}

type postTimelineRow struct {
	PostID      string    `db:"post_id"`
	PublishedAt time.Time `db:"published_at"`
}

func (p *postTimelineRow) toDomain() domain.PostTimeline {
	return domain.PostTimeline{
		PostID:      p.PostID,
		PublishedAt: p.PublishedAt,
	}
}

type TimelineFilter struct {
	Size *int
}
