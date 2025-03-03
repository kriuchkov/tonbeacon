package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-faster/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
	walletutils "github.com/xssnick/tonutils-go/ton/wallet"
)

type MasterKey struct {
	Seed    string              `mapstructure:"seed"`
	Version walletutils.Version `mapstructure:"version"`
}

func (mk *MasterKey) GetSeed() []string {
	return strings.Split(mk.Seed, " ")
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
	LogLevel string         `mapstructure:"log_level"`
	GRPCPort string         `mapstructure:"grpc_port"`
	Master   MasterKey      `mapstructure:"master"`
	Database DatabaseConfig `mapstructure:"database"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName(".config.api")
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
