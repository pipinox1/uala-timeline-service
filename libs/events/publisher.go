package events

import (
	"context"
)

type Publishable interface {
	Key() string
	Topic() string
	Payload() []byte
}

type Publisher interface {
	Publish(ctx context.Context, event Publishable) error
}
