package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"os/signal"

	"github.com/go-faster/errors"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"golang.org/x/sync/errgroup"

	"github.com/kriuchkov/tonbeacon/adapters/consumer"
	"github.com/kriuchkov/tonbeacon/adapters/producer"
	"github.com/kriuchkov/tonbeacon/adapters/repository"
	"github.com/kriuchkov/tonbeacon/ports/outbox"
	"github.com/kriuchkov/tonbeacon/ports/transaction"
)

var (
	enableOutboxProcessor = flag.Bool("outbox-processor", false, "Enable the outbox processor")
	enableKafkaProcessor  = flag.Bool("kafka-processor", false, "Enable the Kafka processor")
)

// Main initializes the processor application, loads configuration, and starts
// enabled processors (Outbox or Kafka) based on flags. It handles graceful
// shutdown via OS signals and logs application activity.
func main() {
	flag.Parse()
	var err error

	// Setup signal handling for graceful shutdown.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	defer log.Info().Msg("processor stopped")

	log.Info().Msg("starting processor")

	var cfg *Config
	if cfg, err = LoadConfig(); err != nil {
		log.Warn().Err(err).Msg("load config")
		os.Exit(64)
	}

	log.Info().Bool("outbox-processor", *enableOutboxProcessor).Bool("kafka-processor", *enableKafkaProcessor).Msg("config loaded")

	if err := cfg.Database.Validate(); err != nil {
		log.Warn().Err(err).Msg("database config validation")
		os.Exit(64)
	}

	dbConnection := pgdriver.NewConnector(pgdriver.WithDSN(cfg.Database.DSN()))

	db := bun.NewDB(sql.OpenDB(dbConnection), pgdialect.New())
	if err := db.PingContext(ctx); err != nil {
		panic(errors.Wrap(err, "ping database").Error())
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Warn().Err(err).Msg("close database")
		}
	}()

	eg := errgroup.Group{}
	// Start the outbox processor if enabled.
	if *enableOutboxProcessor {
		if err := cfg.OutboxProcessor.Validate(); err != nil {
			log.Warn().Err(err).Msg("outbox processor config validation")
			os.Exit(1)
		}

		log.Info().Msg("outbox processor is enabled")

		kafkaProducer, err := producer.NewKafkaProducer(&producer.ProducerOptions{
			Brokers: cfg.OutboxProcessor.Brokers,
			Topic:   cfg.OutboxProcessor.Topic,
			ReqAcks: cfg.OutboxProcessor.RequiredAcks,
			Retries: cfg.OutboxProcessor.MaxRetries,
		})
		if err != nil {
			panic(err.Error())
		}

		outboxConsumer, err := setupOutboxProcessor(db, kafkaProducer)
		if err != nil {
			panic(err.Error())
		}
		defer kafkaProducer.Close() //nolint:errcheck

		eg.Go(func() error { log.Info().Msg("outbox processor started"); outboxConsumer.Consumer(ctx); return nil })
	}

	if *enableKafkaProcessor {
		log.Info().Msg("kafka consumer is enabled")

		txProcessor, err := setupTransactionProcessor(ctx, cfg, db)
		if err != nil {
			panic(err.Error())
		}

		eg.Go(func() error { log.Info().Msg("kafka processor started"); txProcessor.Consume(ctx); return nil })
	}

	eg.Wait()
}

func setupOutboxProcessor(db *bun.DB, writer consumer.OutboxWriter) (*consumer.Outbox, error) {
	outbox := consumer.NewOutbox(consumer.OutboxOptions{
		OutboxManager: outbox.New(repository.New(db)),
		TxManager:     repository.NewTxRepository(db),
		Writer:        writer,
	})
	return outbox, nil
}

func setupTransactionProcessor(ctx context.Context, cfg *Config, db *bun.DB) (*consumer.Kafka, error) {
	dataBase := repository.New(db)
	handler := transaction.New(ctx, &transaction.Options{
		DatabasePort:    dataBase,
		TransactionPort: dataBase,
		TxPort:          repository.NewTxRepository(db),
	})

	kafkaConsumer := consumer.NewKafka(consumer.KafkaOptions{
		Brokers: cfg.TransactionProcessor.Brokers,
		Topic:   cfg.TransactionProcessor.Topic,
		GroupID: cfg.TransactionProcessor.GroupID,
		Handler: handler,
	})
	return kafkaConsumer, nil
}
