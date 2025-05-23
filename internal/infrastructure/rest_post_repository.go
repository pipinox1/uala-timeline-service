package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
	"uala-timeline-service/internal/domain/posts"
)

var _ posts.PostRepository = (*RestPostRepository)(nil)

type RestPostRepository struct {
	client  *resty.Client
	baseURL string
}

func NewRestPostRepository(baseURL string) *RestPostRepository {
	return &RestPostRepository{
		client:  resty.New(),
		baseURL: baseURL,
	}
}

func (r *RestPostRepository) GetPostById(ctx context.Context, id string) (*posts.Post, error) {
	endpoint := fmt.Sprintf("%s/api/v1/posts/%s", r.baseURL, id)

	resp, err := r.client.R().
		SetContext(ctx).
		Get(endpoint)

	if err != nil {
		log.Err(err).Msg("error getting post")
		return nil, fmt.Errorf("error fetching post: %w", err)
	}

	if resp.IsError() {
		log.Err(err).Msg("error getting post")
		return nil, fmt.Errorf("API returned error status: %d - %s", resp.StatusCode(), resp.String())
	}

	var response postResponse
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	return response.toDomain(), nil
}

type multiGetResponse struct {
	Posts []postResponse `json:"posts"`
}

func (r *RestPostRepository) MGetPosts(ctx context.Context, postIDs []string) ([]posts.Post, error) {
	if len(postIDs) == 0 {
		return []posts.Post{}, nil
	}
	idsParam := strings.Join(postIDs, ",")
	endpoint := fmt.Sprintf("%s/api/v1/posts?ids=%s", r.baseURL, idsParam)

	resp, err := r.client.R().
		SetContext(ctx).
		Get(endpoint)

	if err != nil {
		log.Err(err).Msg("error getting mpost")
		return nil, fmt.Errorf("error fetching posts: %w", err)
	}

	if resp.IsError() {
		log.Err(err).Msg("error getting mpost")
		return nil, fmt.Errorf("API returned error status: %d - %s", resp.StatusCode(), resp.String())
	}

	var response multiGetResponse
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	posts := make([]posts.Post, len(response.Posts))
	for i, apiPost := range response.Posts {
		posts[i] = *apiPost.toDomain()
	}

	return posts, nil
}

type postResponse struct {
	ID          string        `json:"id"`
	Contents    []PostContent `json:"contents"`
	AuthorID    string        `json:"author_id"`
	PublishedAt time.Time     `json:"published_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type PostContent struct {
	Type string  `json:"type"`
	Text *string `json:"text,omitempty"`
}

func (p postResponse) toDomain() *posts.Post {
	contents := make([]posts.Content, len(p.Contents))
	for i, content := range p.Contents {
		contents[i] = posts.Content{
			Type: content.Type,
			Text: content.Text,
		}
	}
	return &posts.Post{
		ID:          p.ID,
		Contents:    contents,
		AuthorID:    p.AuthorID,
		PublishedAt: p.PublishedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
