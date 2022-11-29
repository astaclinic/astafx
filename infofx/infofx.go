package infofx

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

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
			p.Logger.Infof("%s\t[INFO]\tStarting application: \n", time.Now().Format(time.RFC3339))
			p.Logger.Infof("%s\t[INFO]\t Golang Version: %v \n", time.Now().Format(time.RFC3339), runtime.Version())
			p.Logger.Infof("%s\t[INFO]\t       Platform: %v %v \n", time.Now().Format(time.RFC3339), runtime.GOOS, runtime.GOARCH)
			hostname, err := os.Hostname()
			if err != nil {
				hostname = fmt.Sprintf("Fail to get hostname: %v", err)
			}
			p.Logger.Infof("%s\t[INFO]\t       Hostname: %v\n", time.Now().Format(time.RFC3339), hostname)
			p.Logger.Infof("%s\t[INFO]\t     Build Date: %v\n", time.Now().Format(time.RFC3339), os.Getenv("BUILD_DATE"))
			p.Logger.Infof("%s\t[INFO]\t   Build Commit: %v\n", time.Now().Format(time.RFC3339), os.Getenv("BUILD_COMMIT"))
			return nil
		},
	})
}

var Module = fx.Module("info",
	fx.Invoke(DisplayInfo),
)
