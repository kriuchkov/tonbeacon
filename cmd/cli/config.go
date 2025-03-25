package main

import (
	"os"
	"strings"

	"github.com/go-faster/errors"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"

	walletutils "github.com/xssnick/tonutils-go/ton/wallet"
)

type MasterKey struct {
	Seed    string              `mapstructure:"seed" validate:"required"`
	Version walletutils.Version `mapstructure:"version" validate:"required"`
}

func (mk *MasterKey) Validate() error {
	if err := validator.New().Struct(mk); err != nil {
		return err
	}
	return nil
}

func (mk *MasterKey) GetSeed() []string {
	return strings.Split(mk.Seed, " ")
}

type Config struct {
	LogLevel string    `mapstructure:"log_level"`
	Master   MasterKey `mapstructure:"master"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName(".config.cli")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME")
	v.AddConfigPath("./.dev")

	v.BindEnv("log_level")
	v.BindEnv("master.seed")

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
