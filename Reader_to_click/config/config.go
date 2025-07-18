package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Clickhouse
		Nats
		Reader
	}

	Clickhouse struct {
		Host     string
		Port     int
		User     string
		Password string
		Dbname   string
	}

	Nats struct {
		URL          string
		ChannelName  string
		ConsumerName string
	}

	Reader struct {
		NatsUrl string
		Streams map[string]StreamConf
	}

	StreamConf struct {
		ChannelName  string
		ConsumerName string
		BatchSize    int
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
