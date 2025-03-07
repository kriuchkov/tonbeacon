package publisher

import (
	"context"

	"github.com/kriuchkov/tonbeacon/core/ports"
)

var _ ports.PublisherPort = (*NoopPublisher)(nil)

type NoopPublisher struct{}

func (p *NoopPublisher) Publish(_ context.Context, _ interface{}) error {
	return nil
}

func (p *NoopPublisher) Close() error {
	return nil
}
