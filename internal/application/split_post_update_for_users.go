package application

import (
	"context"
	"sync"
	"uala-timeline-service/internal/domain"
	"uala-timeline-service/libs/events"
)

type SplitPostUpdateForUsersCommand struct {
	ID       string `json:"id"`
	AuthorID string `json:"author_id"`
}

type SplitPostUpdateForUsers struct {
	postRepository    domain.PostRepository
	followsRepository domain.FollowRepository
	eventPublisher    events.Publisher
}

func NewSplitPostUpdateForUsers(
	postRepository domain.PostRepository,
	followsRepository domain.FollowRepository,
	eventPublisher events.Publisher,
) *SplitPostUpdateForUsers {
	return &SplitPostUpdateForUsers{
		postRepository:    postRepository,
		followsRepository: followsRepository,
		eventPublisher:    eventPublisher,
	}
}

func (s *SplitPostUpdateForUsers) Exec(ctx context.Context, cmd *SplitPostUpdateForUsersCommand) error {
	followerIDs, err := s.followsRepository.GetUserFollowerIDs(ctx, cmd.AuthorID)
	if err != nil {
		return err
	}

	// We do this with goroutines because we will use a best effort approach
	wg := sync.WaitGroup{}
	wg.Add(len(followerIDs))
	for _, followerID := range followerIDs {
		go func(followerID string) {
			defer wg.Done()
			err := s.eventPublisher.Publish(context.WithoutCancel(ctx), domain.NewUserTimelineAddPostEvent(followerID, cmd.ID))
			if err != nil {
				// TODO: log error and send to a retries queue to avoid retrying all the users for some fails
				return
			}
		}(followerID)
	}
	wg.Wait()
	return nil
}
