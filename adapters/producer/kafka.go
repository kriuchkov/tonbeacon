package producer

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/go-faster/errors"
)

type ProducerConfig struct {
	Brokers []string
	Topic   string
}

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(config ProducerConfig) (*KafkaProducer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5
	saramaConfig.Producer.Retry.Backoff = 100 * time.Millisecond

	producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating kafka producer")
	}

	return &KafkaProducer{producer: producer, topic: config.Topic}, nil
}

func (p *KafkaProducer) SendMessage(key string, value []byte) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err = p.producer.SendMessage(msg)
	if err != nil {
		return 0, 0, errors.Wrap(err, "sending message")
	}
	return partition, offset, nil
}

func (p *KafkaProducer) Close() error {
	if err := p.producer.Close(); err != nil {
		return errors.Wrap(err, "closing producer")
	}
	return nil
}
