// Scanner is the entry point for the TON blockchain scanner application.
// It initializes a context with signal handling for graceful shutdown,
// loads application configuration, and establishes a connection to the TON network.
// The application creates a scanner with multiple workers that processes blockchain data
// and publishes the results through a configured publisher.
package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"

	"github.com/kriuchkov/tonbeacon/adapters/publisher"
	"github.com/kriuchkov/tonbeacon/adapters/ton"
	"github.com/kriuchkov/tonbeacon/core/ports"
	"github.com/rs/zerolog/log"
	liteclientutils "github.com/xssnick/tonutils-go/liteclient"
	tonutils "github.com/xssnick/tonutils-go/ton"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	cfg, err := LoadConfig()
	if err != nil {
		log.Panic().Err(err).Msg("config loading")
	}

	client := liteclientutils.NewConnectionPool()
	if err := client.AddConnectionsFromConfigUrl(ctx, cfg.Ton.URL); err != nil {
		log.Panic().Err(err).Msg("liteclient connection")
	}

	liteClient := tonutils.NewAPIClient(client, tonutils.ProofCheckPolicyFast).WithRetry()
	scanner := ton.NewScanner(liteClient, &ton.OptionsScanner{
		NumWorkers: cfg.Scanning.NumWorkers,
	})

	publisher, err := setPublisher(cfg)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create publisher")
	}
	defer publisher.Close()

	resultsCh := make(chan any, 1000)
	if err := scanner.RunAsync(ctx, resultsCh); err != nil {
		log.Panic().Err(err).Msg("run scanner")
	}

	if cfg.PPROF {
		go func() { http.ListenAndServe(":6060", nil) }()
	}

	log.Info().Msg("scanner started")
	for {
		select {
		case result := <-resultsCh:
			if err := publisher.Publish(ctx, result); err != nil {
				log.Error().Err(err).Msg("failed to publish message")
			}
		case <-ctx.Done():
			return
		}
	}
}

// setPublisher creates and returns a publisher based on the provided configuration.
// It supports different publisher types including stdout and Kafka publishers.
//
// Parameters:
//   - cfg: Config containing publisher type and configuration options
//
// Returns:
//   - ports.PublisherPort: An initialized publisher implementation
//   - error: Any error encountered during publisher creation
//
// The publisher type is determined by cfg.PublisherType:
//   - StdoutPublisherType: Returns a stdout publisher
//   - KafkaPublisherType: Returns a Kafka publisher configured with the provided options
//   - Default: Returns a no-operation publisher
func setPublisher(cfg *Config) (ports.PublisherPort, error) {
	switch cfg.PublisherType {
	case StdoutPublisherType:
		return &publisher.StdoutPublisher{}, nil
	case KafkaPublisherType:
		return publisher.NewKafkaPublisher(&publisher.KafkaOptions{
			Brokers:      cfg.Kafka.Brokers,
			Topic:        cfg.Kafka.Topic,
			RequiredAcks: cfg.Kafka.RequiredAcks,
			MaxRetries:   cfg.Kafka.MaxRetries,
		})
	default:
		return &publisher.NoopPublisher{}, nil
	}
}
