package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

type KafkaHandler interface {
	Handle(ctx context.Context, message []byte) error
}

type KafkaConsumerOptions struct {
	Brokers      []string
	Topic        string
	GroupID      string
	Handler      KafkaHandler
	SaramaConfig *sarama.Config
}

type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	handler  KafkaHandler
	topic    string
}

func NewKafkaConsumer(opts KafkaConsumerOptions) *KafkaConsumer {
	consumer, err := sarama.NewConsumerGroup(opts.Brokers, opts.GroupID, opts.SaramaConfig)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create Kafka consumer")
	}

	return &KafkaConsumer{
		consumer: consumer,
		handler:  opts.Handler,
		topic:    opts.Topic,
	}
}

func (c *KafkaConsumer) Consume(ctx context.Context) {
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

func (c *KafkaConsumer) Close() error {
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
			Str("topic", message.Topic).
			Int32("partition", message.Partition).
			Int64("offset", message.Offset).
			Msg("received message")

		if err := h.handler.Handle(h.ctx, message.Value); err != nil {
			log.Error().Err(err).
				Str("topic", message.Topic).
				Int32("partition", message.Partition).
				Int64("offset", message.Offset).
				Msg("failed to handle message")
		} else {
			// Только если успешно обработали, отмечаем как обработанное
			session.MarkMessage(message, "")
		}
	}
	return nil
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
