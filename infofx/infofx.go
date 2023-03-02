package infofx

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

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
			p.Logger.Infof(" Golang Version: %v", runtime.Version())
			p.Logger.Infof("       Platform: %v %v", runtime.GOOS, runtime.GOARCH)
			hostname, err := os.Hostname()
			if err != nil {
				hostname = fmt.Sprintf("Fail to get hostname: %v", err)
			}
			p.Logger.Infof("       Hostname: %v", hostname)
			buildInfo, ok := debug.ReadBuildInfo()
			if !ok {
				return fmt.Errorf("failed to read build info")
			}
			buildCommit := os.Getenv("BUILD_COMMIT")
			for _, buildSetting := range buildInfo.Settings {
				if buildSetting.Key == "vcs.revision" {
					buildCommit = buildSetting.Value
				}
			}
			p.Logger.Infof("     Build Date: %v", BuildDate)
			p.Logger.Infof("   Build Commit: %v", buildCommit)
			return nil
		},
	})
}

var Module = fx.Module("info",
	fx.Invoke(DisplayInfo),
)
