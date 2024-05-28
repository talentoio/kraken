package config

import (
	"context"
	"time"

	configLib "kraken/pkg/config"
)

type Config struct {
	Server         ServerConfig `validate:"required"`
	Kraken         KrakenConfig `validate:"required"`
	Cache          CacheConfig  `validate:"required"`
	AvailablePairs []string     `validate:"required,min=1"`
}

type (
	ServerConfig struct {
		Port int `validate:"required,min=80"`
	}

	KrakenConfig struct {
		HOST           string        `validate:"required"`
		UpdateDuration time.Duration `validate:"required"`
		RequestTimeOut time.Duration `validate:"required"`
	}

	CacheConfig struct {
		Expire time.Duration `validate:"required"`
	}
)

func GetConfig(ctx context.Context) (*Config, error) {

	var cfg Config
	err := configLib.Parse(
		ctx,
		configLib.Options{
			Dir:  "./config",
			File: "config.yaml",
			Type: "yaml",
		},
		&cfg,
	)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
