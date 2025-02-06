package config

import (
	env "github.com/caarlos0/env/v11"
	"github.com/palomachain/paloma-cdp/internal/app/gql"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
)

type Config struct {
	Persistence persistence.Configuration
	GraphQL     gql.Configuration
}

func Parse() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
