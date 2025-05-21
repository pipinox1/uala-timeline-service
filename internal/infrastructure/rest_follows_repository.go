package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"uala-timeline-service/internal/domain"
)

var _ domain.FollowRepository = (*RestFollowsRepository)(nil)

type RestFollowsRepository struct {
	client  *resty.Client
	baseURL string
}

func (r *RestFollowsRepository) GetUserFollowerIDs(ctx context.Context, userID string) ([]string, error) {
	endpoint := fmt.Sprintf("%s/api/v1/follow/user/%s/followers", r.baseURL, userID)
	resp, err := r.client.R().
		SetContext(ctx).
		Get(endpoint)

	if err != nil {
		return nil, fmt.Errorf("error fetching post: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API returned error status: %d - %s", resp.StatusCode(), resp.String())
	}

	var response followersResponse
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	return response.toDomain(), nil
}

func NewRestFollowsRepository(baseURL string) *RestFollowsRepository {
	return &RestFollowsRepository{
		client:  resty.New(),
		baseURL: baseURL,
	}
}

type followersResponse struct {
	Followers []string `json:"followers"`
}

func (p followersResponse) toDomain() []string {
	return p.Followers
}
