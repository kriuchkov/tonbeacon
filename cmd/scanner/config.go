//nolint:errcheck
package main

import (
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/go-faster/errors"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
)

const (
	// defaultLogLevel is the default log level.
	defaultLogLevel = "info"

	// defaultTestnetConfigURL is the default URL for the testnet config.
	defaultTestnetConfigURL = "https://tonutils.com/testnet-global.config.json"

	// defaultKafkaMaxRetries is the default number of retries for Kafka producer.
	defaultKafkaMaxRetries = 3

	// defaultKafkaRequiredAcks is the default number of required acks for Kafka producer.
	defaultKafkaRequiredAcks = sarama.WaitForAll

	// defaultScanningNumWorkers is the default number of workers for scanning.
	defaultScanningNumWorkers = 1

	// defaultPublisherType is the default publisher type.
	defaultPublisherType = StdoutPublisherType
)

type PublisherType string

const (
	NoopPublisherType   PublisherType = "none"
	StdoutPublisherType PublisherType = "stdout"
	KafkaPublisherType  PublisherType = "kafka"
)

type KafkaConfig struct {
	Brokers      []string            `mapstructure:"brokers"`
	Topic        string              `mapstructure:"topic"`
	MaxRetries   int                 `mapstructure:"max_retries"`
	RequiredAcks sarama.RequiredAcks `mapstructure:"required_acks"`
}

type ScanningConfig struct {
	NumWorkers int `mapstructure:"num_workers" validate:"required"`
}

type TonConfig struct {
	URL string `mapstructure:"url" validate:"required"`
}

type Config struct {
	LogLevel string      `mapstructure:"log_level"`
	PPROF    string      `mapstructure:"pprof"`
	Kafka    KafkaConfig `mapstructure:"kafka"`

	// required
	PublisherType PublisherType  `mapstructure:"publisher_type" validate:"required"`
	Scanning      ScanningConfig `mapstructure:"scanning"  validate:"required"`
	Ton           TonConfig      `mapstructure:"ton" validate:"required"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	// file
	v.SetConfigName(".config.scanner2")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME")

	// env
	v.SetEnvPrefix("tonbeacon")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables
	v.BindEnv("log_level")
	v.BindEnv("pprof")
	v.BindEnv("kafka.brokers")
	v.BindEnv("kafka.topic")
	v.BindEnv("kafka.max_retries")
	v.BindEnv("kafka.required_acks")
	v.BindEnv("publisher_type")
	v.BindEnv("scanning.num_workers")
	v.BindEnv("ton.url")

	v.SetDefault("log_level", defaultLogLevel)
	v.SetDefault("kafka.max_retries", defaultKafkaMaxRetries)
	v.SetDefault("kafka.required_acks", defaultKafkaRequiredAcks)
	v.SetDefault("scanning.num_workers", defaultScanningNumWorkers)
	v.SetDefault("ton.url", defaultTestnetConfigURL)
	v.SetDefault("publisher_type", defaultPublisherType)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errors.Wrap(err, "read config")
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, errors.Wrap(err, "unmarshal config")
	}

	level, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, errors.Wrap(err, "parse log level")
	}

	zerolog.SetGlobalLevel(level)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"

	if err := validator.New().Struct(&config); err != nil {
		return nil, errors.Wrap(err, "validate config")
	}
	return &config, nil
}
