package cqrx

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/cyg-pd/go-kebabcase"
)

type eventProcessorOption interface {
	apply(*cqrs.EventProcessorConfig)
}
type eventProcessorOptionFunc func(*cqrs.EventProcessorConfig)

func (fn eventProcessorOptionFunc) apply(cfg *cqrs.EventProcessorConfig) { fn(cfg) }

func WithEventProcessorOnHandle(fn func(params cqrs.EventProcessorOnHandleParams) error) eventProcessorOption {
	return eventProcessorOptionFunc(func(c *cqrs.EventProcessorConfig) { c.OnHandle = fn })
}

func WithEventProcessorMarshaler(m cqrs.CommandEventMarshaler) eventProcessorOption {
	return eventProcessorOptionFunc(func(c *cqrs.EventProcessorConfig) { c.Marshaler = m })
}

func WithEventProcessorGenerateSubscribeTopic(fn cqrs.EventProcessorGenerateSubscribeTopicFn) eventProcessorOption {
	return eventProcessorOptionFunc(func(c *cqrs.EventProcessorConfig) { c.GenerateSubscribeTopic = fn })
}

func defaultEventProcessorGenerateSubscribeTopic(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
	return "events-" + kebabcase.Kebabcase(params.EventName), nil
}

func defaultEventProcessorConfig() cqrs.EventProcessorConfig {
	return cqrs.EventProcessorConfig{
		GenerateSubscribeTopic: defaultEventProcessorGenerateSubscribeTopic,
		Marshaler:              DefaultMarshaler(),
		Logger:                 watermill.NewSlogLogger(nil),
	}
}
