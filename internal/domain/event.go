package domain

import "encoding/json"

const (
	UserTimelineAddPostTopic = "user_timeline.add_post"
)

type UserTimelineAddPostEvent struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func (p UserTimelineAddPostEvent) Key() string {
	return p.ID
}

func (p UserTimelineAddPostEvent) Topic() string {
	return UserTimelineAddPostTopic
}

func (p UserTimelineAddPostEvent) Payload() []byte {
	payload, _ := json.Marshal(p)
	return payload
}

func NewUserTimelineAddPostEvent(userID string, postID string) UserTimelineAddPostEvent {
	return UserTimelineAddPostEvent{ID: postID, UserID: userID}
}
