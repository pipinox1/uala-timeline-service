package domain

import "encoding/json"

const (
	UserTimelineAddPostTopic = "user_timeline.add_post"
)

type UserTimelineAddPostEvent struct {
	PostID string `json:"post_id"`
	UserID string `json:"user_id"`
}

func (p UserTimelineAddPostEvent) Key() string {
	return p.PostID
}

func (p UserTimelineAddPostEvent) Topic() string {
	return UserTimelineAddPostTopic
}

func (p UserTimelineAddPostEvent) Payload() []byte {
	payload, _ := json.Marshal(p)
	return payload
}

func NewUserTimelineAddPostEvent(userID string, postID string) UserTimelineAddPostEvent {
	return UserTimelineAddPostEvent{PostID: postID, UserID: userID}
}
