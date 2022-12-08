package sentryfx

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/fx"
)

type SentryConfig struct {
	Dsn string `mapstructure:"dsn" yaml:"dsn"  validate:"required,uri"`
}

func RunSentry(lifecycle fx.Lifecycle, config *SentryConfig) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return sentry.Init(sentry.ClientOptions{
				Dsn: config.Dsn,
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
