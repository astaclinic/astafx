package redisfx

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Module("redis",
	fx.Provide(New),
)

type RedisConfig struct {
	Dsn      string `mapstructure:"dsn" yaml:"dsn" validate:"required,hostname_port"`
	Password string `mapstructure:"password" yaml:"password" validate:"printascii"`
}

func init() {
	// config must have a default value for viper to load config from env variables
	// default value of empty string (zero value) will not pass the "required" config validation
	viper.SetDefault("redis.dsn", "")
	viper.SetDefault("redis.password", "")
}

func New(config *RedisConfig) (*redis.Client, error) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     config.Dsn,
		Password: config.Password, // no password set
		DB:       0,               // use default DB
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
