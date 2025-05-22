package events

import (
	"context"
)

type Publishable interface {
	Key() string
	Topic() string
	Payload() []byte
}

//go:generate mockery --name=Publisher --output=mocks --outpkg=mocks_events
type Publisher interface {
	Publish(ctx context.Context, event Publishable) error
}
