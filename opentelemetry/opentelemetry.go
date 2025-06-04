package opentelemetry

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/cyg-pd/go-watermillx/internal/utils"
	"github.com/cyg-pd/go-watermillx/opentelemetry/metrics"
	wotelfloss "github.com/dentech-floss/watermill-opentelemetry-go-extra/pkg/opentelemetry"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
)

func PublisherDecorator(metricBuilder metrics.Builder) message.PublisherDecorator {
	return func(pub message.Publisher) (message.Publisher, error) {
		name := utils.StructName(pub)

		var err error
		pub, err = metricBuilder.DecoratePublisher(pub)
		if err != nil {
			return nil, err
		}

		pub = wotel.NewNamedPublisherDecorator(name, pub)
		pub = wotelfloss.NewTracePropagatingPublisherDecorator(pub)
		return pub, nil
	}
}

func SubscriberDecorator(metricBuilder metrics.Builder) message.SubscriberDecorator {
	return func(sub message.Subscriber) (message.Subscriber, error) {
		return metricBuilder.DecorateSubscriber(sub)
	}
}

func HandlerMiddleware(metricBuilder metrics.Builder) message.HandlerMiddleware {
	return func(h message.HandlerFunc) message.HandlerFunc {
		h = metricBuilder.RouterMiddleware(h)
		h = wotelfloss.ExtractRemoteParentSpanContext()(h)
		h = wotel.Trace()(h)
		return h
	}
}
