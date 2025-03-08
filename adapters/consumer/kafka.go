package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

type KafkaHandler interface {
	Handle(ctx context.Context, message []byte) error
}

type KafkaOptions struct {
	Brokers []string
	Topic   string
	GroupID string
	Handler KafkaHandler
}

type Kafka struct {
	consumer sarama.ConsumerGroup
	handler  KafkaHandler
	topic    string
}

func NewKafka(opts KafkaOptions) *Kafka {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(opts.Brokers, opts.GroupID, cfg)
	if err != nil {
		panic(err.Error())
	}

	return &Kafka{consumer: consumer, handler: opts.Handler, topic: opts.Topic}
}

func (c *Kafka) Consume(ctx context.Context) {
	go func() {
		for err := range c.consumer.Errors() {
			log.Error().Err(err).Msg("kafka consumer error")
		}
	}()

	handler := &consumerGroupHandler{ctx: ctx, handler: c.handler}
	topics := []string{c.topic}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := c.consumer.Consume(ctx, topics, handler); err != nil {
				log.Error().Err(err).Msg("error from consumer")
			}
		}
	}
}

func (c *Kafka) Close() error {
	return c.consumer.Close()
}

type consumerGroupHandler struct {
	ctx     context.Context
	handler KafkaHandler
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if h.ctx.Err() != nil {
			return nil //nolint:nilerr
		}

		log.Debug().
			Str("topic", message.Topic).Int32("partition", message.Partition).Int64("offset", message.Offset).
			Msg("received message")

		if err := h.handler.Handle(h.ctx, message.Value); err != nil {
			log.Error().Err(err).
				Str("topic", message.Topic).Int32("partition", message.Partition).Int64("offset", message.Offset).
				Msg("handle message")
		} else {
			session.MarkMessage(message, "")
		}
	}
	return nil
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

type StdOutHandler struct{}

func (h *StdOutHandler) Handle(_ context.Context, message []byte) error {
	log.Info().Bytes("message", message).Msg("received message")
	return nil
}
