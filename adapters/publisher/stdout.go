package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/kriuchkov/tonbeacon/core/ports"
)

var _ ports.PublisherPort = (*StdoutPublisher)(nil)

type StdoutPublisher struct{}

func (p *StdoutPublisher) Publish(_ context.Context, message any) error {
	data, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, string(data))
	return nil
}

func (p *StdoutPublisher) Close() error {
	return nil
}
