package driver

import (
	"fmt"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/cyg-pd/go-reflectx"
	"github.com/cyg-pd/go-watermillx/metrics"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
)

var store sync.Map

func Register(name string, driver DriverFunc) {
	store.Store(name, driver)
}

func New(name string, config any) (Driver, error) {
	if d, ok := store.Load(name); ok {
		return newOTELWrapper(d.(DriverFunc)(config))
	}

	return nil, fmt.Errorf("watermillx/driver: %s driver not found", name)
}

type otelWrapper struct {
	Next Driver
}

// Publisher implements Driver.
func (w *otelWrapper) Publisher() (message.Publisher, error) {
	pub, err := w.Next.Publisher()
	if err != nil {
		return nil, err
	}

	pub, err = metrics.Builder.DecoratePublisher(pub)
	if err != nil {
		return nil, fmt.Errorf("watermillx/driver: could not apply publisher otel metric decorator: %w", err)
	}

	return wotel.NewNamedPublisherDecorator(reflectx.PackageName(), pub), nil
}

// Subscriber implements Driver.
func (w *otelWrapper) Subscriber(topic string, opts ...SubscriberOption) (message.Subscriber, error) {
	sub, err := w.Next.Subscriber(topic, opts...)
	if err != nil {
		return nil, err
	}

	sub, err = metrics.Builder.DecorateSubscriber(sub)
	if err != nil {
		return nil, fmt.Errorf("watermillx/driver: could not apply subscriber otel metric decorator: %w", err)
	}

	return sub, nil
}

func newOTELWrapper(driver Driver, err error) (Driver, error) {
	if err != nil {
		return nil, err
	}
	return &otelWrapper{Next: driver}, nil
}
