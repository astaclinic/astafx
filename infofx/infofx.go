package infofx

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"go.uber.org/fx"
	"go.uber.org/zap"
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
			p.Logger.Infof(" Golang Version: %v", runtime.Version())
			p.Logger.Infof("       Platform: %v %v", runtime.GOOS, runtime.GOARCH)
			hostname, err := os.Hostname()
			if err != nil {
				hostname = fmt.Sprintf("Fail to get hostname: %v", err)
			}
			p.Logger.Infof("       Hostname: %v", hostname)
			p.Logger.Infof("     Build Date: %v", os.Getenv("BUILD_DATE"))
			p.Logger.Infof("   Build Commit: %v", os.Getenv("BUILD_COMMIT"))
			return nil
		},
	})
}

var Module = fx.Module("info",
	fx.Invoke(DisplayInfo),
)
