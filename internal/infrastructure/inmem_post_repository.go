package infrastructure

import (
	"context"
	"uala-timeline-service/internal/domain"
	"uala-timeline-service/utils"
)

var _ domain.PostRepository = (*InmemPostRepository)(nil)

type InmemPostRepository struct {
	posts map[string]domain.Post
}

func NewInmemPostRepository() *InmemPostRepository {
	return &InmemPostRepository{
		posts: make(map[string]domain.Post),
	}
}

func (r *InmemPostRepository) MGetPosts(ctx context.Context, postIDs []string) ([]domain.Post, error) {
	mockPosts := make([]domain.Post, len(postIDs))
	for i, postID := range postIDs {
		mockPosts[i] = domain.Post{
			ID: postID,
			Contents: []domain.Content{
				domain.Content{
					Type: "text",
					Text: utils.Ref("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur"),
				},
			},
		}
	}
	return mockPosts, nil
}

func (r *InmemPostRepository) GetPostById(ctx context.Context, id string) (*domain.Post, error) {
	return &domain.Post{
		ID: id,
		Contents: []domain.Content{
			domain.Content{
				Type: "text",
				Text: utils.Ref("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur"),
			},
		},
	}, nil
}
