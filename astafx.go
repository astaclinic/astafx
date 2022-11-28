package astafx

import (
	"github.com/astaclinic/astafx/httpfx"
	"github.com/astaclinic/astafx/loggerfx"
	"github.com/astaclinic/astafx/metricsfx"
	"github.com/astaclinic/astafx/routerfx"
	"github.com/astaclinic/astafx/sentryfx"
	"go.uber.org/fx"
)

var Module = fx.Module("asta",
	httpfx.Module,
	loggerfx.Module,
	metricsfx.Module,
	routerfx.Module,
	sentryfx.Module,
)
