package routerfx

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("router",
	fx.Provide(New),
)

type Params struct {
	fx.In
	Logger           *zap.SugaredLogger `optional:"true"`
	ControllerRoutes []ControllerRoute  `group:"controllerRoutes"`
	HandlerRoutes    []HandlerRoute     `group:"handlerRoutes"`
}

type Result struct {
	fx.Out
	Http *gin.Engine
}

func New(p Params) Result {
	gin.SetMode(gin.ReleaseMode)

	http := gin.New()
	if p.Logger != nil {
		http.Use(ginzap.Ginzap(p.Logger.Desugar(), time.RFC3339, true))
	}
	http.Use(gin.Recovery())

	apiRouterGroup := http.Group("/v1")
	for _, route := range p.ControllerRoutes {
		if p.Logger != nil {
			p.Logger.Infow("registering controller route", "pattern", route.RoutePattern())
		}
		route.RegisterControllerRoutes(
			apiRouterGroup.Group(route.RoutePattern()),
		)
	}

	for _, route := range p.HandlerRoutes {
		if p.Logger != nil {
			p.Logger.Infow("registering handler route", "pattern", route.RoutePattern())
		}
		http.Any(route.RoutePattern(), gin.WrapH(route.HttpHandler()))
	}

	return Result{
		Http: http,
	}
}

func (r *Result) GetHttpRouter() *gin.Engine {
	return r.Http
}

type ControllerRoute interface {
	RegisterControllerRoutes(rg *gin.RouterGroup)
	RoutePattern() string
}

func AsControllerRoute(controller any) any {
	return fx.Annotate(
		controller,
		fx.As(new(ControllerRoute)),
		fx.ResultTags(`group:"controllerRoutes"`),
	)
}

type HandlerRoute interface {
	HttpHandler() http.Handler
	RoutePattern() string
}

func AsHandlerRoute(handler any) any {
	return fx.Annotate(
		handler,
		fx.As(new(HandlerRoute)),
		fx.ResultTags(`group:"handlerRoutes"`),
	)
}
