package httpfx

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Module("http",
	fx.Provide(NewHttp),
	fx.Invoke(RunHttpServer),
)

type HttpConfig struct {
	ListenAddr string `mapstructure:"listen_addr" required:"required,hostname_port"`
}

type HttpParams struct {
	fx.In
	Config  *HttpConfig
	Handler *gin.Engine
}

func NewHttp(p HttpParams) *http.Server {
	return &http.Server{
		Addr:    p.Config.ListenAddr,
		Handler: p.Handler,
	}
}

type RunHttpParams struct {
	fx.In
	Lifecycle  fx.Lifecycle
	HttpServer *http.Server
}

func RunHttpServer(p RunHttpParams) {
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go p.HttpServer.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return p.HttpServer.Shutdown(ctx)
		},
	})
}
