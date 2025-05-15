package domain

import "context"

type FollowRepository interface {
	GetUserFollowerIDs(ctx context.Context, userID string) ([]string, error)
}
