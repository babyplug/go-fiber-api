package config

import (
	"go-fiber-api/internal/wrapper/logx"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	_configOnce sync.Once
	// All values below are default if they are not specified from env variable
	_config = Configuration{
		DevMode:       false,
		Port:          "80",
		IsAutoMigrate: true,
		TZ:            "Asia/Bangkok",
	}

	_log *logrus.Entry
)

type Configuration struct {
	DevMode bool   `mapstructure:"DEV_MODE"`
	Port    string `mapstructure:"PORT"`
	TZ      string `mapstructure:"TZ"`

	DBUser        string `mapstructure:"DB_USER"`
	DBPass        string `mapstructure:"DB_PASS"`
	DBName        string `mapstructure:"DB_NAME"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBSSLMode     string `mapstructure:"DB_SSL_MODE"`
	IsAutoMigrate bool   `mapstructure:"IS_AUTO_MIGRATE"`

	// Redis
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	// API Keys
	SecretKey string `mapstructure:"SECRET_KEY"`

	// CORS
	CorsAllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	CorsAllowedHeaders string `mapstructure:"CORS_ALLOWED_HEADERS"`
}

func Provide(logx *logx.LogX) *Configuration {
	_configOnce.Do(func() {
		_log = logx.WithFields(logrus.Fields{
			"component": "config",
			"module":    "config",
		})

		envFilePath, ok := os.LookupEnv("ENV_FILE_PATH")
		if len(envFilePath) > 0 && ok {
			viper.SetConfigFile(envFilePath)
			if err := viper.ReadInConfig(); err != nil {
				_log.Warnf("read config from file %s failed: %v\tcontinue reading from `env`\n", envFilePath, err)
			}
		} else {
			viper.AutomaticEnv()
		}

		bindEnv(_config)

		err := viper.Unmarshal(&_config)
		if err != nil {
			logrus.Fatalf("config bind failed: %v", err)
		}

		_log.Infof("Configuration: %v", logx.GetLevel())
		_log.Debugf("Configuration: %+v", _config)
	})

	return &_config
}

func ResetProvide() {
	_configOnce = sync.Once{}
}

func bindEnv(dest interface{}, parts ...string) {
	ifv := reflect.ValueOf(dest)
	ift := reflect.TypeOf(dest)

	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)

		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}

		switch v.Kind() {
		case reflect.Struct:
			bindEnv(v.Interface(), append(parts, tv)...)
		default:
			envKey := strings.Join(append(parts, tv), ".")
			err := viper.BindEnv(envKey)
			if err != nil {
				_log.Printf("bind env key %s failed: %v\n", envKey, err)
			}
		}
	}
}
