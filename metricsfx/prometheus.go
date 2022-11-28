package metricsfx

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusHandler struct{}

func NewPrometheusHandler() *PrometheusHandler {
	return &PrometheusHandler{}
}

func (ph *PrometheusHandler) HttpHandler() http.Handler {
	return promhttp.Handler()
}

func (ph *PrometheusHandler) RoutePattern() string {
	return "/metrics"
}
