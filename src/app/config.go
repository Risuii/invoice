package app

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/spf13/viper"
)

/*
	All config should be required.
	Optional only allowed if zero value of the type is expected being the default value.
	time.Duration units are “ns”, “us” (or “µs”), “ms”, “s”, “m”, “h”. as in time.ParseDuration().
*/

type (
	Postgres struct {
		ConnURI            string        `mapstructure:"PG_CONN_URI" validate:"required"`
		MaxPoolSize        int           `mapstructure:"PG_MAX_POOL_SZE"` //Optional, default to 0 (zero value of int)
		MaxIdleConnections int           `mapstructure:"PG_MAX_IDLE_CONNECTIONS"`
		MaxIdleTime        time.Duration `mapstructure:"PG_MAX_IDLE_TIME"` //Optional, default to '0s' (zero value of time.Duration)
		MaxLifeTime        time.Duration `mapstructure:"PG_MAX_IDLE_TIME"` //Optional, default to '0s' (zero value of time.Duration)
	}

	Redis struct {
		Host     string `mapstructure:"REDIS_HOST" validate:"required"`
		Password string `mapstructure:"REDIS_PASSWORD"`
	}

	Configuration struct {
		ServiceName string      `mapstructure:"SERVICE_NAME"`
		Postgres    Postgres    `mapstructure:",squash"`
		Redis       Redis       `mapstructure:",squash"`
		Translation Translation `mapstructure:",squash"`

		Environment string `mapstructure:"ENV" validate:"required,oneof=development staging production"`
		BindAddress int    `mapstructure:"BIND_ADDRESS" validate:"required"`
		LogLevel    int    `mapstructure:"LOG_LEVEL" validate:"required"`
	}
)

func InitConfig(ctx context.Context) (*Configuration, error) {
	var cfg Configuration

	viper.SetConfigType("env")
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	_, err := os.Stat(envFile)
	if !os.IsNotExist(err) {
		viper.SetConfigFile(envFile)

		if err := viper.ReadInConfig(); err != nil {
			log.Printf("failed to read config:%v", err)
			return nil, err
		}
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("failed to bind config:%v", err)
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Printf("invalid config:%v", err)
		}
		log.Println("failed to load config")
		return nil, err
	}

	log.Printf("Config loaded: %+v", cfg)
	return &cfg, nil
}
