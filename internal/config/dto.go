package config

import (
	"go.uber.org/zap/zapcore"
)

type configDTO struct {
	Logger   loggerConfigDTO `envPrefix:"LOG_" json:"log"`
	API      apiConfigDTO    `envPrefix:"API_" json:"api"`
	Database dbConfigDTO     `envPrefix:"DB_" json:"db"`
}

type loggerConfigDTO struct {
	Level zapcore.Level `env:"LEVEL" envDefault:"debug" json:"level"`
}

type apiConfigDTO struct {
	Host string `env:"HOST" envDefault:"0.0.0.0" json:"host"`
	Port int    `env:"PORT" envDefault:"8080" json:"port"`
}

type dbConfigDTO struct {
	Host     string `env:"HOST" envDefault:"localhost" json:"host"`
	Port     int    `env:"PORT" envDefault:"5432" json:"port"`
	Username string `env:"USERNAME" envDefault:"postgres" json:"username"`
	Password string `env:"PASSWORD,unset" envDefault:"postgres" json:"password"`
	Database string `env:"DATABASE" envDefault:"postgres" json:"database"`
}
