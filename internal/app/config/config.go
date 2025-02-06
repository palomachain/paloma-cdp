package config

import (
	env "github.com/caarlos0/env/v11"
)

func Parse[T any]() (cfg *T, err error) {
	return cfg, env.Parse(cfg)
}
