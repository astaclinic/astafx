package metricsfx

import (
	"github.com/astaclinic/astafx/routerfx"
	"go.uber.org/fx"
)

var Module = fx.Module("metrics",
	fx.Provide(routerfx.AsHandlerRoute(NewPrometheusHandler)),
)
