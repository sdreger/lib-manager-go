package config

import "github.com/caarlos0/env/v11"

const envPrefix = "LIB_MANAGER_"

func New() (AppConfig, error) {
	var cfg AppConfig
	if err := env.ParseWithOptions(&cfg, env.Options{Prefix: envPrefix}); err != nil {
		return cfg, err
	}

	return cfg, nil
}
