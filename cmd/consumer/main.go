package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/kriuchkov/tonbeacon/adapters/consumer"
	"github.com/kriuchkov/tonbeacon/adapters/repository"
	"github.com/kriuchkov/tonbeacon/ports/outbox"
	"github.com/kriuchkov/tonbeacon/ports/transaction"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	cfg, err := LoadConfig()
	if err != nil {
		log.Panic().Err(err).Msg("config loading")
	}

	db := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.Database.DSN()))), pgdialect.New())

	if err := db.PingContext(ctx); err != nil {
		log.Panic().Err(err).Msg("db connection")
	}

	outboxManager := outbox.New(repository.New(db))

	outboxConsumer := consumer.NewOutbox(consumer.OutboxOptions{
		OutboxManager: outboxManager,
		TxManager:     nil,
		Writer:        nil,
	})

	log.Info().Msg("starting outbox consumer")
	go outboxConsumer.Consumer(ctx)

	handler := transaction.New(ctx, &transaction.Options{})

	kafkaConsumer := consumer.NewKafkaConsumer(consumer.KafkaConsumerOptions{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
		GroupID: cfg.Kafka.GroupID,
		Handler: handler,
		//SaramaConfig: config,
	})

	log.Info().Msg("starting kafka consumer")
	go kafkaConsumer.Consume(ctx)

	<-ctx.Done()
	log.Info().Msg("consumer stopped")
}
