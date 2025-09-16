package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

func New() (Config, error) {
	var dto configDTO

	if err := env.Parse(&dto); err != nil {
		return Config{}, fmt.Errorf("parse config from the env: %w", err)
	}

	return newConfig(dto), nil
}

func newConfig(dto configDTO) Config {
	conf := Config{
		log: LoggerConfig{
			level: dto.Logger.Level,
		},
		api: APIConfig{
			host: dto.API.Host,
			port: dto.API.Port,
		},
		db: DbConfig{
			host:     dto.Database.Host,
			port:     dto.Database.Port,
			username: dto.Database.Username,
			password: dto.Database.Password,
			database: dto.Database.Database,
		},
	}

	maskPasswordVars(
		&dto.Database.Password,
	)

	conf.origin = dto

	return conf
}
