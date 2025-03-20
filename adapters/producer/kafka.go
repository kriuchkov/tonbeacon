package producer

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/go-faster/errors"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
)

const (
	// defaultBackoff is the default backoff duration for the producer.
	defaultBackoff = 100 * time.Millisecond

	// defaultRetries is the default number of retries for the producer.
	defaultRetries = 5

	// defaultReqAcks is the default number of required acks for the producer.
	defaultReqAcks = sarama.WaitForAll

	// defaultSucceed is the default value for the producer to return successes.
	defaultSucceed = true

	// defaultError is the default value for the producer to return errors.
	defaultError = true
)

type ProducerOptions struct {
	Brokers []string `required:"true"`
	Topic   string   `required:"true"`
	Succeed *bool
	Error   *bool
	ReqAcks sarama.RequiredAcks
	Retries int
	Backoff time.Duration
}

func (c *ProducerOptions) SetDefaults() {
	if c.Succeed == nil {
		c.Succeed = lo.ToPtr(defaultSucceed)
	}

	if c.Error == nil {
		c.Error = lo.ToPtr(defaultError)
	}
	if c.ReqAcks == 0 {
		c.ReqAcks = defaultReqAcks
	}
	if c.Retries == 0 {
		c.Retries = defaultRetries
	}
	if c.Backoff == 0 {
		c.Backoff = defaultBackoff
	}
}

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(opt *ProducerOptions) (*KafkaProducer, error) {
	opt.SetDefaults()

	if err := validator.New().Struct(opt); err != nil {
		return nil, errors.Wrap(err, "validating producer options")
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = *opt.Succeed
	saramaConfig.Producer.Return.Errors = *opt.Error
	saramaConfig.Producer.RequiredAcks = opt.ReqAcks
	saramaConfig.Producer.Retry.Max = opt.Retries
	saramaConfig.Producer.Retry.Backoff = opt.Backoff

	producer, err := sarama.NewSyncProducer(opt.Brokers, saramaConfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating kafka producer")
	}

	return &KafkaProducer{producer: producer, topic: opt.Topic}, nil
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
