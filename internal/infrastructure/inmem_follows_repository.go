package infrastructure

import (
	"context"
	"github.com/google/uuid"
	"uala-timeline-service/internal/domain"
)

var _ domain.FollowRepository = (*InmemFollowsRepository)(nil)

type InmemFollowsRepository struct {
}

func (i InmemFollowsRepository) GetUserFollowerIDs(ctx context.Context, userID string) ([]string, error) {
	return []string{
		uuid.New().String(),
		uuid.New().String(),
	}, nil
}

func NewInmemFollowsRepository() *InmemFollowsRepository {
	return &InmemFollowsRepository{}
}
