package main

import (
	"os"

	"github.com/IBM/sarama"
	"github.com/go-faster/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
)

const (
	defautlTestnetConfigURL = "https://tonutils.com/testnet-global.config.json"
)

type PublisherType string

const (
	NoopPublisherType   PublisherType = ""
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
	NumWorkers int `mapstructure:"num_workers"`
}

type TonConfig struct {
	URL string `mapstructure:"url"`
}

type Config struct {
	LogLevel      string
	PublisherType PublisherType  `mapstructure:"publisher_type"`
	Kafka         KafkaConfig    `mapstructure:"kafka"`
	Scanning      ScanningConfig `mapstructure:"scanning"`
	Ton           TonConfig      `mapstructure:"ton"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName(".config.scanner")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME")

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "read config")
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, errors.Wrap(err, "unmarshal config")
	}

	if config.Ton.URL == "" {
		config.Ton.URL = defautlTestnetConfigURL
	}

	level, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, errors.Wrap(err, "parse log level")
	}

	zerolog.SetGlobalLevel(level)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	return &config, nil
}
