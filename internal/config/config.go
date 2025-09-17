package config

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"go.uber.org/zap/zapcore"
)

// Config is an application configuration read-only container.
type Config struct {
	log LoggerConfig
	api APIConfig
	db  DbConfig

	origin configDTO
}

func (c Config) Logger() LoggerConfig {
	return c.log
}

func (c Config) API() APIConfig {
	return c.api
}

func (c Config) Db() DbConfig {
	return c.db
}

func (c Config) DebugJSON() json.RawMessage {
	data, err := json.Marshal(c.origin)
	if err == nil {
		return data
	}

	fallback := map[string]string{
		"config_debug":   fmt.Sprintf("%+v", c.origin),
		"marshall_error": err.Error(),
	}

	//nolint:errchkjson
	data, _ = json.Marshal(fallback)

	return data
}

type LoggerConfig struct {
	level zapcore.Level
}

func (c LoggerConfig) Level() zapcore.Level {
	return c.level
}

type APIConfig struct {
	host string
	port int
}

func (c APIConfig) Host() string {
	return c.host
}

func (c APIConfig) Port() int {
	return c.port
}

func (c APIConfig) Address() string {
	return net.JoinHostPort(c.Host(), strconv.Itoa(c.Port()))
}

type DbConfig struct {
	host     string
	port     int
	username string
	password string
	database string
	sslMode  string
}

func (c DbConfig) Host() string {
	return c.host
}

func (c DbConfig) Port() int {
	return c.port
}

func (c DbConfig) Address() string {
	return net.JoinHostPort(c.Host(), strconv.Itoa(c.Port()))
}

func (c DbConfig) Username() string {
	return c.username
}

func (c DbConfig) Password() string {
	return c.password
}

func (c DbConfig) Database() string {
	return c.database
}

func (c DbConfig) ToDSN() string {
	return c.toDSN(false)
}

func (c DbConfig) ToDSNDebug() string {
	return c.toDSN(true)
}

func (c DbConfig) toDSN(maskPwd bool) string {
	password := c.Password()
	if maskPwd {
		password = maskPassword(password)
	}

	sslMode := c.sslMode
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dsnVal(c.Username()), dsnVal(password), c.Host(), c.Port(), dsnVal(c.Database()), sslMode,
	)
}

func dsnVal(value string) string {
	if !strings.Contains(value, " ") {
		return value
	}

	return "'" + strings.ReplaceAll(value, "'", `\'`) + "'"
}
