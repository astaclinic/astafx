package redisfx

import (
	"context"

	"github.com/go-redis/redis/v9"
	"go.uber.org/fx"
)

var Module = fx.Module("redis",
	fx.Provide(New),
)

type RedisConfig struct {
	Dsn      string `mapstructure:"dsn" yaml:"dsn"  validate:"required,hostname"`
	Password string `mapstructure:"password" yaml:"password"  validate:"required,printascii"`
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
