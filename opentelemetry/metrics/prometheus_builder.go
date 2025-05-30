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

// AttachRouter implements MetricsBuilder.
func (p *PrometheusMetricsBuilder) AttachRouter(router *message.Router) {
	p.AddPrometheusRouterMetrics(router)
}

func NewPrometheusMetricsBuilder(prometheusRegistry prometheus.Registerer, namespace string, subsystem string) *PrometheusMetricsBuilder {
	return &PrometheusMetricsBuilder{
		PrometheusMetricsBuilder: metrics.NewPrometheusMetricsBuilder(prometheusRegistry, namespace, subsystem),
	}
}

var _ Builder = (*PrometheusMetricsBuilder)(nil)
