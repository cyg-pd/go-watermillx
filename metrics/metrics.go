package metrics

import (
	"github.com/ThreeDotsLabs/watermill/components/metrics"
	"github.com/cyg-pd/go-reflectx"
	"go.opentelemetry.io/otel"
)

var Builder = metrics.NewOpenTelemetryMetricsBuilder(otel.Meter(reflectx.PackageName()), "", "")
