package metrics

import (
	"github.com/ThreeDotsLabs/watermill/components/metrics"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusMetricsBuilder provides methods to decorate publishers, subscribers and handlers.
type PrometheusMetricsBuilder struct {
	metrics.PrometheusMetricsBuilder
}

// RouterMiddleware implements Builder.
func (p *PrometheusMetricsBuilder) RouterMiddleware(h message.HandlerFunc) message.HandlerFunc {
	return p.NewRouterMiddleware().Middleware(h)
}
func NewPrometheusMetricsBuilder(prometheusRegistry prometheus.Registerer, namespace string, subsystem string) *PrometheusMetricsBuilder {
	return &PrometheusMetricsBuilder{
		PrometheusMetricsBuilder: metrics.NewPrometheusMetricsBuilder(prometheusRegistry, namespace, subsystem),
	}
}

var _ Builder = (*PrometheusMetricsBuilder)(nil)
