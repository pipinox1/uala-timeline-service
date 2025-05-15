package domain

import (
	"context"
)

type PostRepository interface {
	MGetPosts(ctx context.Context, postIDs []string) ([]Post, error)
	GetPostById(ctx context.Context, id string) (*Post, error)
}

type Post struct {
	ID       string
	Contents []Content
	AuthorID string
}

type Content struct {
	Type string
	Text *string
	Url  *string
}
