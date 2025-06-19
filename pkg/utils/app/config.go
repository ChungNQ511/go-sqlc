package app

import (
	"time"

	"github.com/spf13/viper"
)

var AppConfiguration AppConfig

//nolint:revive // This disables revive for this function

type AppConfig struct {
	// Environment
	ENVIRONMENT       string `mapstructure:"GO_ENV"`
	SERVER_ADDRESS    string `mapstructure:"SERVER_ADDRESS"`
	CORS_ORIGIN       string `mapstructure:"CORS_ORIGIN"`
	SERVICE_HOST_NAME string `mapstructure:"SERVICE_HOST_NAME"`

	// Database
	DATABASE_DRIVER   string `mapstructure:"DATABASE_DRIVER"`
	DATABASE_HOST     string `mapstructure:"DATABASE_HOST"`
	DATABASE_PORT     string `mapstructure:"DATABASE_PORT"`
	DATABASE_USERNAME string `mapstructure:"DATABASE_USERNAME"`
	DATABASE_PASSWORD string `mapstructure:"DATABASE_PASSWORD"`
	DATABASE_NAME     string `mapstructure:"DATABASE_NAME"`

	// Redis
	REDIS_HOST         string        `mapstructure:"REDIS_HOST"`
	REDIS_PASSWORD     string        `mapstructure:"REDIS_PASSWORD"`
	REDIS_MAX_IDLE     int           `mapstructure:"REDIS_MAX_IDLE"`
	REDIS_MAX_ACTIVE   int           `mapstructure:"REDIS_MAX_ACTIVE"`
	REDIS_IDLE_TIMEOUT time.Duration `mapstructure:"REDIS_IDLE_TIMEOUT"`

	// SWAGGER
	SWAGGER_HOST string `mapstructure:"SWAGGER_HOST"`

	// S3
	BUCKET            string `mapstructure:"BUCKET"`
	ACCESS_KEY_ID     string `mapstructure:"ACCESS_KEY_ID"`
	ACCESS_KEY_SECRET string `mapstructure:"ACCESS_KEY_SECRET"`
	S3_REGION         string `mapstructure:"S3_REGION"`
	S3_HOST_NAME      string `mapstructure:"S3_HOST_NAME"`

	// APM
	ELASTIC_APM_SERVER_URL   string `mapstructure:"ELASTIC_APM_SERVER_URL"`
	ELASTIC_APM_SERVICE_NAME string `mapstructure:"ELASTIC_APM_SERVICE_NAME"`
	ELASTIC_APM_ENVIRONMENT  string `mapstructure:"ELASTIC_APM_ENVIRONMENT"`
	ELASTIC_APM_SECRET_TOKEN string `mapstructure:"ELASTIC_APM_SECRET_TOKEN"`

	// CLIENTS
	KEYCLOAK_URL                 string `mapstructure:"KEYCLOAK_URL"`
	KEYCLOAK_CLIENT_ID           string `mapstructure:"KEYCLOAK_CLIENT_ID"`
	KEYCLOAK_CLIENT_SECRET       string `mapstructure:"KEYCLOAK_CLIENT_SECRET"`
	KEYCLOAK_AUDIENCE            string `mapstructure:"KEYCLOAK_AUDIENCE"`
	KEYCLOAK_SCOPE               string `mapstructure:"KEYCLOAK_SCOPE"`
	KEYCLOAK_TOKEN_CACHE_SECONDS int    `mapstructure:"KEYCLOAK_TOKEN_CACHE_SECONDS"`
}

// LoadAppConfig load the app config
func LoadAppConfig(path string) (config AppConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
