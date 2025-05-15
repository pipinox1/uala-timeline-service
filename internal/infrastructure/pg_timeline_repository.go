package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"time"
	"uala-timeline-service/internal/domain"
)

var _ domain.TimelineRepository = (*TimelineRepository)(nil)

type TimelineRepository struct {
	db *sqlx.DB
}

func (t TimelineRepository) AddPostToTimeline(ctx context.Context, userID string, postID string) error {
	//TODO implement me
	panic("implement me")
}

func NewTimelineRepository(db *sqlx.DB) *TimelineRepository {
	return &TimelineRepository{db: db}
}

func (t TimelineRepository) GetUserTimeline(ctx context.Context, userID string, filter domain.TimelineFilter) (*domain.UserTimeline, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()

	sb.Select("user_id", "post_id", "created_at")
	sb.From("timelines")
	sb.Where(sb.Equal("user_id", userID))
	sb.OrderBy("created_at DESC")

	if filter.Size != nil {
		sb.Limit(*filter.Size)
	}

	query, args := sb.Build()

	var rows []timelineRow
	err := t.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting user timeline: %w", err)
	}

	var postIDs []string
	var lastUpdate time.Time
	for _, row := range rows {
		postIDs = append(postIDs, row.PostID)
		if row.AddedAt.After(lastUpdate) {
			lastUpdate = row.AddedAt
		}
	}

	if len(postIDs) == 0 {
		return &domain.UserTimeline{
			UserID:     userID,
			Posts:      []string{},
			LastUpdate: time.Time{},
		}, nil
	}

	return &domain.UserTimeline{
		UserID:     userID,
		Posts:      postIDs,
		LastUpdate: lastUpdate,
	}, nil
}

func (t TimelineRepository) SaveTimeline(ctx context.Context, timeline *domain.UserTimeline) error {
	tx, err := t.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	stmt, err := tx.PreparexContext(ctx, `
		INSERT INTO timelines (user_id, post_id, created_at)
		VALUES ($1, $2, $3)`)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, postID := range timeline.Posts {
		_, err = stmt.ExecContext(ctx, timeline.UserID, postID, timeline.LastUpdate)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error inserting timeline entry: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

type timelineRow struct {
	UserID  string    `db:"user_id"`
	PostID  string    `db:"post_id"`
	AddedAt time.Time `db:"created_at"`
}

type TimelineFilter struct {
	Size *int
}
