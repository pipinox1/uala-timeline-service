package application

import (
	"time"
	"uala-timeline-service/internal/domain"
)

type TimelineFilled struct {
	LastUpdate time.Time `json:"last_update"`
	Posts      []Post    `json:"posts"`
	UserID     string    `json:"user_id"`
}

type Post struct {
	ID       string    `json:"id"`
	Contents []Content `json:"contents"`
	AuthorID string    `json:"author_id"`
}

type Content struct {
	Type string  `json:"type"`
	Text *string `json:"text"`
	Url  *string `json:"url"`
}

func FromDomain(timelineFilled *domain.TimelineFilled) *TimelineFilled {
	posts := make([]Post, len(timelineFilled.Posts))
	for i, post := range timelineFilled.Posts {
		contents := make([]Content, len(post.Contents))
		for i, content := range post.Contents {
			contents[i] = Content{
				Type: content.Type,
				Text: content.Text,
				Url:  content.Url,
			}
		}
		posts[i] = Post{
			ID:       post.ID,
			Contents: contents,
			AuthorID: post.AuthorID,
		}
	}

	return &TimelineFilled{
		LastUpdate: timelineFilled.LastUpdate,
		Posts:      posts,
		UserID:     timelineFilled.UserID,
	}
}
