package trace

import (
	"github.com/ThreeDotsLabs/watermill/message"
	wotelfloss "github.com/dentech-floss/watermill-opentelemetry-go-extra/pkg/opentelemetry"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
)

func PublisherDecorator() message.PublisherDecorator {
	return func(pub message.Publisher) (message.Publisher, error) {
		pub = wotelfloss.NewTracePropagatingPublisherDecorator(pub)
		pub = wotel.NewPublisherDecorator(pub)
		return pub, nil
	}
}

func HandlerMiddleware() message.HandlerMiddleware {
	return func(h message.HandlerFunc) message.HandlerFunc {
		h = wotelfloss.ExtractRemoteParentSpanContext()(h)
		h = wotel.Trace()(h)
		return h
	}
}
