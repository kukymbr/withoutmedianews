package config_test

import (
	"fmt"
	"testing"

	"github.com/kukymbr/withoutmedianews/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	type testData struct {
		name   string
		env    map[string]string
		assert func(t *testing.T, conf config.Config, err error)
	}

	tests := []testData{
		{
			name: "when values are set",
			env: map[string]string{
				"LOG_LEVEL":   "info",
				"API_PORT":    "8181",
				"API_HOST":    "127.0.0.1",
				"DB_HOST":     "database",
				"DB_PORT":     "5431",
				"DB_USERNAME": "test",
				"DB_PASSWORD": "test",
				"DB_DATABASE": "test",
			},
			assert: func(t *testing.T, conf config.Config, err error) {
				require.NoError(t, err)

				assert.Equal(t, zapcore.InfoLevel, conf.Logger().Level())
				assert.Equal(t, "127.0.0.1:8181", conf.API().Address())
				assert.Equal(t, "postgres://test:test@database:5431/test?sslmode=disable", conf.Db().ToDSN())
				assert.Equal(t, "postgres://test:*****@database:5431/test?sslmode=disable", conf.Db().ToDSNDebug())
			},
		},
		{
			name: "when values are not set",
			env:  map[string]string{},
			assert: func(t *testing.T, conf config.Config, err error) {
				require.NoError(t, err)

				assert.Equal(t, zapcore.DebugLevel, conf.Logger().Level())
				assert.Equal(t, "0.0.0.0:8080", conf.API().Address())
				assert.Equal(t, "localhost", conf.Db().Host())
				assert.Equal(t, 5432, conf.Db().Port())
				assert.Equal(t, "postgres", conf.Db().Username())
				assert.Equal(t, "postgres", conf.Db().Password())
				assert.Equal(t, "postgres", conf.Db().Database())
			},
		},
	}

	invalidCases := map[string]string{
		"LOG_LEVEL": "unknown",
		"API_PORT":  "NaN",
		"DB_PORT":   "NaN",
	}

	for key, val := range invalidCases {
		tests = append(tests, testData{
			name: "when invalid: " + key,
			env: map[string]string{
				key: val,
			},
			assert: func(t *testing.T, conf config.Config, err error) {
				assert.Error(t, err)
			},
		})
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("from env: %s", test.name), func(t *testing.T) {
			for k, v := range test.env {
				t.Setenv(k, v)
			}

			conf, err := config.New()
			test.assert(t, conf, err)
		})
	}
}

func TestConfig_DebugJSON(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)

	assert.NotEmpty(t, conf.DebugJSON())
}
