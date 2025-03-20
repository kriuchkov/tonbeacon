package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/go-faster/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
)

const (
	// defaultEnableConsumer is the default value for EnableConsumer.
	defaultEnableConsumer = false

	// defaultLogLevel is the default log level.
	defaultLogLevel = "info"

	// defaultKafkaMaxRetries is the default number of retries for Kafka producer.
	defaultKafkaMaxRetries = 3

	// defaultKafkaRequiredAcks is the default number of required acks for Kafka producer.
	defaultKafkaRequiredAcks = sarama.WaitForAll
)

type KafkaConfig struct {
	Brokers      []string            `mapstructure:"brokers"`
	Topic        string              `mapstructure:"topic"`
	GroupID      string              `mapstructure:"group_id"`
	MaxRetries   int                 `mapstructure:"max_retries"`
	RequiredAcks sarama.RequiredAcks `mapstructure:"required_acks"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode" default:"disable"`
}

func (dc *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", dc.User, dc.Password, dc.Host, dc.Port, dc.DBName, dc.SSLMode)
}

type Config struct {
	// Enables consumers
	EnableOutboxConsumer bool `mapstructure:"enable_outbox_consumer"`
	EnableKafkaConsumer  bool `mapstructure:"enable_kafka_consumer"`

	PPROF    string `mapstructure:"pprof"`
	LogLevel string `mapstructure:"log_level"`

	// Configs for consumers
	Database DatabaseConfig `mapstructure:"database"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	// file
	v.SetConfigName(".config.consumer")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME")

	// env
	v.SetEnvPrefix("tonbeacon")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables
	v.BindEnv("enable_outbox_consumer")
	v.BindEnv("enable_kafka_consumer")
	v.BindEnv("log_level")
	v.BindEnv("database.host")
	v.BindEnv("database.port")
	v.BindEnv("database.user")
	v.BindEnv("database.password")
	v.BindEnv("database.dbname")
	v.BindEnv("database.sslmode")
	v.BindEnv("kafka.brokers")
	v.BindEnv("kafka.topic")
	v.BindEnv("kafka.group_id")
	v.BindEnv("kafka.max_retries")
	v.BindEnv("kafka.required_acks")

	// Defaults
	v.SetDefault("enable_outbox_consumer", defaultEnableConsumer)
	v.SetDefault("enable_kafka_consumer", defaultEnableConsumer)
	v.SetDefault("log_level", defaultLogLevel)
	v.SetDefault("kafka.max_retries", defaultKafkaMaxRetries)
	v.SetDefault("kafka.required_acks", defaultKafkaRequiredAcks)

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
