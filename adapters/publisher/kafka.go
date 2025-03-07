package publisher

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-faster/errors"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"

	"github.com/kriuchkov/tonbeacon/core/ports"
)

var _ ports.PublisherPort = (*KafkaPublisher)(nil)

type KafkaOptions struct {
	Brokers      []string `required:"true"`
	Topic        string   `required:"true"`
	RequiredAcks sarama.RequiredAcks
	MaxRetries   int
}

func (k *KafkaOptions) SetDefaults() {
	if k.RequiredAcks == 0 {
		k.RequiredAcks = sarama.WaitForAll
	}
	if k.MaxRetries == 0 {
		k.MaxRetries = 3
	}
}

type KafkaPublisher struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaPublisher(opt *KafkaOptions) (*KafkaPublisher, error) {
	if err := validator.New().Struct(opt); err != nil {
		log.Panic().Err(err).Msg("kafka options")
	}

	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = opt.RequiredAcks
	cfg.Producer.Retry.Max = opt.MaxRetries
	cfg.Producer.Return.Successes = true

	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrap(err, "kafka config")
	}

	producer, err := sarama.NewSyncProducer(opt.Brokers, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "kafka producer")
	}
	return &KafkaPublisher{producer: producer, topic: opt.Topic}, nil
}

func (p *KafkaPublisher) Publish(ctx context.Context, message any) error {
	data, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "json marshal")
	}

	msg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Value:     sarama.StringEncoder(data),
		Timestamp: time.Now(),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return errors.Wrap(err, "send message to Kafka")
	}

	log.Debug().
		Str("topic", p.topic).Int32("partition", partition).Int64("offset", offset).
		Msg("message sent to Kafka")
	return nil
}

func (p *KafkaPublisher) Close() error {
	return p.producer.Close()
}
