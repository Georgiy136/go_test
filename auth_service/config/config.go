package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Http               `yaml:"http"`
		Postgres           `yaml:"postgres"`
		NotificationClient `yaml:"notification_client"`
		Tokens
		AccessToken  `yaml:"accesstoken"`
		RefreshToken `yaml:"refreshtoken"`
		Crypter      `yaml:"crypter"`
	}

	Http struct {
		Port int
	}

	Postgres struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
		Sslmode  string `yaml:"sslmode"`
	}

	AccessToken struct {
		SignedKey     string `yaml:"signed_key"`
		TokenLifetime string `yaml:"token_lifetime"`
	}
	RefreshToken struct {
		SignedKey     string `yaml:"signed_key"`
		TokenLifetime string `yaml:"token_lifetime"`
	}
	Crypter struct {
		SignedKey string `yaml:"signed_key"`
	}
	Tokens struct {
		AccessToken
		RefreshToken
	}

	NotificationClient struct {
		Url string
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

	logrus.Infof("config: %+v", cfg)

	cfg.Tokens = Tokens{
		AccessToken: AccessToken{
			SignedKey:     "abcdabcd",
			TokenLifetime: "2h",
		},
		RefreshToken: RefreshToken{
			SignedKey:     "abcdabcd",
			TokenLifetime: "1s",
		},
	}
	cfg.Crypter = Crypter{
		SignedKey: "abcdabcdabcdabcd",
	}

	return cfg, nil
}
