package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type MasterKey struct {
	Seed    []string
	Version int
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
	GRPCPort string         `mapstructure:"grpc_port"`
	Master   MasterKey      `mapstructure:"master"`
	Database DatabaseConfig `mapstructure:"database"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	return &Config{}, nil
}
