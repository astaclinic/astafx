package sentryfx

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var (
	Release     string
	Environment string
)

type SentryConfig struct {
	Dsn   string `mapstructure:"dsn" yaml:"dsn" validate:"required,uri"`
	Debug bool   `mapstructure:"debug" yaml:"debug"`
}

func init() {
	// config must have a default value for viper to load config from env variables
	// default value of empty string (zero value) will not pass the "required" config validation
	viper.SetDefault("sentry.dsn", "")
}

func RunSentry(lifecycle fx.Lifecycle, config *SentryConfig) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return sentry.Init(sentry.ClientOptions{
				Dsn:         config.Dsn,
				Release:     Release,
				Environment: Environment,
				Debug:       config.Debug,
			})
		},
		OnStop: func(ctx context.Context) error {
			sentry.Flush(2 * time.Second)
			return nil
		},
	})
}

var Module = fx.Module("sentry",
	fx.Invoke(RunSentry),
)
