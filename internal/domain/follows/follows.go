package follows

import "context"

//go:generate mockery --name=FollowRepository --filename=mocks_follow_repository.go
type FollowRepository interface {
	GetUserFollowerIDs(ctx context.Context, userID string) ([]string, error)
}
