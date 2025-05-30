package cqrx

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/cyg-pd/go-kebabcase"
)

type eventBusOption interface{ apply(*cqrs.EventBusConfig) }
type eventBusOptionFunc func(*cqrs.EventBusConfig)

func (fn eventBusOptionFunc) apply(cfg *cqrs.EventBusConfig) { fn(cfg) }

func WithEventBusOnPublish(fn func(params cqrs.OnEventSendParams) error) eventBusOption {
	return eventBusOptionFunc(func(c *cqrs.EventBusConfig) { c.OnPublish = fn })
}

func WithEventBusMarshaler(m cqrs.CommandEventMarshaler) eventBusOption {
	return eventBusOptionFunc(func(c *cqrs.EventBusConfig) { c.Marshaler = m })
}

func WithEventBusGeneratePublishTopic(fn cqrs.GenerateEventPublishTopicFn) eventBusOption {
	return eventBusOptionFunc(func(c *cqrs.EventBusConfig) { c.GeneratePublishTopic = fn })
}

func defaultEventBusGeneratePublishTopic(params cqrs.GenerateEventPublishTopicParams) (string, error) {
	return "events-" + kebabcase.Kebabcase(params.EventName), nil
}

func defaultEventBusConfig() cqrs.EventBusConfig {
	return cqrs.EventBusConfig{
		GeneratePublishTopic: defaultEventBusGeneratePublishTopic,
		Marshaler:            DefaultMarshaler(),
		Logger:               watermill.NewSlogLogger(nil),
	}
}
