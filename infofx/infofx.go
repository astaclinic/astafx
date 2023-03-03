package infofx

import (
	"context"

	"github.com/astaclinic/astafx/info"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	BuildDate string
)

type InfoParams struct {
	fx.In
	Lifecycle fx.Lifecycle
	Logger    *zap.SugaredLogger
}

func DisplayInfo(p InfoParams) {
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			p.Logger.Info("Starting application:")
			info, err := info.GetInfo()
			if err != nil {
				return err
			}
			p.Logger.Info("\n" + info)
			return nil
		},
	})
}

var Module = fx.Module("info",
	fx.Invoke(DisplayInfo),
)
