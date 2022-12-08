package httpfx

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Module("http",
	fx.Provide(NewHttp),
	fx.Invoke(RunHttpServer),
)

type HttpConfig struct {
	ListenAddr string `mapstructure:"listen_addr" yaml:"listen_addr" required:"required,hostname_port"`
}

func init() {
	// config must have a default value for viper to load config from env variables
	// default value of empty string (zero value) will not pass the "required" config validation
	viper.SetDefault("http.listen_addr", ":8080")
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
			go func() {
				if err := p.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return p.HttpServer.Shutdown(ctx)
		},
	})
}
