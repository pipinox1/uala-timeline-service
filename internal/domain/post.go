package domain

import (
	"context"
	"time"
)

type PostRepository interface {
	MGetPosts(ctx context.Context, postIDs []string) ([]Post, error)
	GetPostById(ctx context.Context, id string) (*Post, error)
}

type Post struct {
	ID          string
	Contents    []Content
	AuthorID    string
	PublishedAt time.Time
	UpdatedAt   time.Time
}

type Content struct {
	Type string
	Text *string
	Url  *string
}
