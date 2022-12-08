package metricsfx

import (
	"go.uber.org/fx"

	"github.com/astaclinic/astafx/routerfx"
)

var Module = fx.Module("metrics",
	fx.Provide(routerfx.AsHandlerRoute(NewPrometheusHandler)),
)
