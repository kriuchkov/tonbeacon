package main

import (
	"fmt"
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

	// defaultKafkaMaxRetries is the default number of retries for Kafka producer.
	defaultKafkaMaxRetries = 3

	// defaultKafkaRequiredAcks is the default number of required acks for Kafka producer.
	defaultKafkaRequiredAcks = sarama.WaitForAll
)

type DatabaseConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     int    `mapstructure:"port" validate:"required"`
	User     string `mapstructure:"user" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	DBName   string `mapstructure:"dbname" validate:"required"`
	SSLMode  string `mapstructure:"sslmode" default:"disable"`
}

func (dc *DatabaseConfig) Validate() error {
	if err := validator.New().Struct(dc); err != nil {
		return errors.Wrap(err, "validate database config")
	}
	return nil
}

func (dc *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", dc.User, dc.Password, dc.Host, dc.Port, dc.DBName, dc.SSLMode)
}

type Kafka struct {
	Brokers      []string            `mapstructure:"brokers" validate:"required"`
	Topic        string              `mapstructure:"topic" validate:"required"`
	GroupID      string              `mapstructure:"group_id"`
	MaxRetries   int                 `mapstructure:"max_retries"`
	RequiredAcks sarama.RequiredAcks `mapstructure:"required_acks"`
}

type TransactionProcessorConfig struct {
	Kafka `mapstructure:",squash"`
}

type OutboxProcessorConfig struct {
	Kafka `mapstructure:",squash" validate:"required"`
}

func (oc *OutboxProcessorConfig) Validate() error {
	if err := validator.New().Struct(oc); err != nil {
		return errors.Wrap(err, "validate database config")
	}
	return nil
}

type Config struct {
	PPROF    string         `mapstructure:"pprof"`
	LogLevel string         `mapstructure:"log_level"`
	Database DatabaseConfig `mapstructure:"database"`

	// Configs for processors
	TransactionProcessor TransactionProcessorConfig `mapstructure:"transaction_processor"`
	OutboxProcessor      OutboxProcessorConfig      `mapstructure:"outbox_processor"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	// file
	v.SetConfigName(".config.consumer")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME")
	v.AddConfigPath("./.dev")

	// env
	v.SetEnvPrefix("tonbeacon")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables
	v.BindEnv("log_level")
	v.BindEnv("database.host")
	v.BindEnv("database.port")
	v.BindEnv("database.user")
	v.BindEnv("database.password")
	v.BindEnv("database.dbname")
	v.BindEnv("database.sslmode")

	// outbox processor
	v.BindEnv("outbox_processor.brokers")
	v.BindEnv("outbox_processor.topic")
	v.BindEnv("outbox_processor.max_retries")
	v.BindEnv("outbox_processor.required_acks")

	// transaction processor
	v.BindEnv("transaction_processor.brokers")
	v.BindEnv("transaction_processor.topic")
	v.BindEnv("transaction_processor.group_id")
	v.BindEnv("transaction_processor.max_retries")
	v.BindEnv("transaction_processor.required_acks")

	// Defaults
	v.SetDefault("log_level", defaultLogLevel)

	v.SetDefault("outbox_processor.max_retries", defaultKafkaMaxRetries)
	v.SetDefault("outbox_processor.required_acks", defaultKafkaRequiredAcks)
	v.SetDefault("transaction_processor.max_retries", defaultKafkaMaxRetries)
	v.SetDefault("transaction_processor.required_acks", defaultKafkaRequiredAcks)

	if err := v.ReadInConfig(); err != nil {
		var errViper viper.ConfigFileNotFoundError
		if !errors.As(err, &errViper) {
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
	return &config, nil
}
