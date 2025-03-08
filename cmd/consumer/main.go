package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"

	"github.com/go-faster/errors"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"golang.org/x/sync/errgroup"

	"github.com/kriuchkov/tonbeacon/adapters/consumer"
	"github.com/kriuchkov/tonbeacon/adapters/repository"
	"github.com/kriuchkov/tonbeacon/ports/outbox"
	//"github.com/kriuchkov/tonbeacon/ports/transaction"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	defer log.Info().Msg("consumer stopped")

	log.Info().Msg("starting consumer")

	cfg, err := LoadConfig()
	if err != nil {
		log.Panic().Err(err).Msg("config loading")
	}

	log.Info().
		Bool("enable_outbox_consumer", cfg.EnableOutboxConsumer).Bool("enable_kafka_consumer", cfg.EnableKafkaConsumer).
		Msg("config loaded")

	eg := errgroup.Group{}
	if cfg.EnableOutboxConsumer {
		log.Info().Msg("outbox consumer is enabled")

		outboxConsumer, err := setupOutboxConsumer(ctx, cfg)
		if err != nil {
			panic(err.Error())
		}

		eg.Go(func() error { log.Info().Msg("starting outbox consumer"); outboxConsumer.Consumer(ctx); return nil })
	}

	if cfg.EnableKafkaConsumer {
		log.Info().Msg("kafka consumer is enabled")

		kafkaConsumer, err := setupKafkaConsumer(ctx, cfg)
		if err != nil {
			panic(err.Error())
		}

		eg.Go(func() error { log.Info().Msg("starting kafka consumer"); kafkaConsumer.Consume(ctx); return nil })
	}

	eg.Wait()
}

func setupOutboxConsumer(ctx context.Context, cfg *Config) (*consumer.Outbox, error) {
	db := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.Database.DSN()))), pgdialect.New())

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "db connection")
	}

	outboxManager := outbox.New(repository.New(db))

	outboxConsumer := consumer.NewOutbox(consumer.OutboxOptions{
		OutboxManager: outboxManager,
		TxManager:     nil,
		Writer:        nil,
	})

	return outboxConsumer, nil
}

func setupKafkaConsumer(ctx context.Context, cfg *Config) (*consumer.Kafka, error) {
	//handler := transaction.New(ctx, &transaction.Options{})

	kafkaConsumer := consumer.NewKafka(consumer.KafkaOptions{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
		GroupID: cfg.Kafka.GroupID,
		Handler: &consumer.StdOutHandler{},
	})

	return kafkaConsumer, nil
}
