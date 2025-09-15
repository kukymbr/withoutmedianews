package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/kukymbr/withoutmedianews/internal/util"
)

func New() (Config, error) {
	var dto configDTO

	if err := env.Parse(&dto); err != nil {
		return Config{}, fmt.Errorf("parse config from the env: %w", err)
	}

	return configFromDTO(dto), nil
}

func configFromDTO(dto configDTO) Config {
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

	util.MaskPasswordVars(
		&dto.Database.Password,
	)

	conf.origin = dto

	return conf
}
