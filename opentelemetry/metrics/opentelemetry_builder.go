package metrics

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/cyg-pd/go-watermillx/internal/utils"
	"github.com/cyg-pd/go-watermillx/messagex"
	"go.opentelemetry.io/otel/metric"
)

func NewOpenTelemetryMetricsBuilder(meter metric.Meter, namespace string, subsystem string) OpenTelemetryMetricsBuilder {
	return OpenTelemetryMetricsBuilder{
		Namespace: namespace,
		Subsystem: subsystem,
		meter:     meter,
	}
}

// OpenTelemetryMetricsBuilder provides methods to decorate publishers, subscribers and handlers.
type OpenTelemetryMetricsBuilder struct {
	meter metric.Meter

	Namespace string
	Subsystem string
	// PublishBuckets defines the histogram buckets for publish time histogram, defaulted if nil.
	PublishBuckets []float64
	// HandlerBuckets defines the histogram buckets for handle execution time histogram, defaulted to watermill's default.
	HandlerBuckets []float64
}

// DecoratePublisher wraps the underlying publisher with OpenTelemetry metrics.
func (b OpenTelemetryMetricsBuilder) DecoratePublisher(pub message.Publisher) (message.Publisher, error) {
	var err error
	d := PublisherOpenTelemetryMetricsDecorator{
		pub:           pub,
		publisherName: utils.StructName(pub),
	}

	d.publishTimeSeconds, err = b.meter.Float64Histogram(
		b.name("publish_time_seconds"),
		metric.WithUnit("seconds"),
		metric.WithDescription("The time that a publishing attempt (success or not) took in seconds"),
		metric.WithExplicitBucketBoundaries(b.PublishBuckets...),
	)

	if err != nil {
		return nil, fmt.Errorf("could not register publish time metric: %w", err)
	}
	return d, nil
}

// DecorateSubscriber wraps the underlying subscriber with OpenTelemetry metrics.
func (b OpenTelemetryMetricsBuilder) DecorateSubscriber(sub message.Subscriber) (message.Subscriber, error) {
	var err error
	d := &SubscriberOpenTelemetryMetricsDecorator{
		subscriberName: utils.StructName(sub),
	}

	d.subscriberMessagesReceivedTotal, err = b.meter.Int64Counter(
		b.name("subscriber_messages_received_total"),
		metric.WithDescription("The total number of messages received by the subscriber"),
	)
	if err != nil {
		return nil, fmt.Errorf("could not register time to ack metric: %w", err)
	}

	d.Subscriber, err = messagex.MessageTransformSubscriberDecorator(d.recordMetrics)(sub)
	if err != nil {
		return nil, fmt.Errorf("could not decorate subscriber with metrics decorator: %w", err)
	}

	return d, nil
}
func (b OpenTelemetryMetricsBuilder) name(name string) string {
	if b.Subsystem != "" {
		name = b.Subsystem + "_" + name
	}
	if b.Namespace != "" {
		name = b.Namespace + "_" + name
	}
	return name
}

var _ Builder = (*OpenTelemetryMetricsBuilder)(nil)
