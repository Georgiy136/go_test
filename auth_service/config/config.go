package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Http
		Postgres
		Tokens
	}

	Http struct {
		Port int
	}

	Postgres struct {
		Host     string
		Port     int
		User     string
		Password string
		Dbname   string
		Sslmode  string
	}

	Tokens struct {
		AccessToken
		RefreshToken
		Crypter
	}
	AccessToken struct {
		SignedKey     string
		tokenLifetime string
	}
	RefreshToken struct {
		SignedKey     string
		tokenLifetime string
	}
	Crypter struct {
		SignedKey string
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("config file read error: %w", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	return cfg, nil
}
