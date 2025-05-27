package router

import (
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/cyg-pd/go-watermillx/metrics"
	wotelfloss "github.com/dentech-floss/watermill-opentelemetry-go-extra/pkg/opentelemetry"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
)

type option interface{ apply(*message.RouterConfig) }
type optionFunc func(*message.RouterConfig)

func (fn optionFunc) apply(cfg *message.RouterConfig) { fn(cfg) }

func WithCloseTimeout(t time.Duration) option {
	return optionFunc(func(cfg *message.RouterConfig) { cfg.CloseTimeout = t })
}

func New(opts ...option) (*message.Router, error) {
	var c message.RouterConfig
	for _, opt := range opts {
		opt.apply(&c)
	}

	r, err := message.NewRouter(c, watermill.NewSlogLogger(nil))
	if r != nil {
		r.AddMiddleware(middleware.Recoverer)
		r.AddMiddleware(wotelfloss.ExtractRemoteParentSpanContext())
		r.AddMiddleware(wotel.Trace())
		metrics.Builder.AddOpenTelemetryRouterMetrics(r)
	}

	return r, err
}
