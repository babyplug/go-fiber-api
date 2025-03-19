package config_test

import (
	"go-fiber-api/internal/core/config"
	"go-fiber-api/internal/wrapper/logx"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	logger := logx.Provide()

	tests := []struct {
		name       string
		expected   config.Configuration
		isClearEnv bool
	}{
		{
			name:       "if_no_env_file_path_should_get_default_config",
			isClearEnv: true,
			expected: config.Configuration{
				DevMode:       false,
				Port:          "80",
				IsAutoMigrate: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.isClearEnv {
				os.Setenv("ENV_FILE_PATH", "")
			}

			// test logic
			cfg := config.Provide(logger)
			assert.Equal(t, test.expected, *cfg)
		})
	}
}
