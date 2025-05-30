package cqrx

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/cyg-pd/go-kebabcase"
)

type eventGroupProcessorOption interface {
	apply(*cqrs.EventGroupProcessorConfig)
}
type eventGroupProcessorOptionFunc func(*cqrs.EventGroupProcessorConfig)

func (fn eventGroupProcessorOptionFunc) apply(cfg *cqrs.EventGroupProcessorConfig) { fn(cfg) }

func WithEventGroupProcessorOnHandle(fn func(params cqrs.EventGroupProcessorOnHandleParams) error) eventGroupProcessorOption {
	return eventGroupProcessorOptionFunc(func(c *cqrs.EventGroupProcessorConfig) { c.OnHandle = fn })
}

func WithEventGroupProcessorMarshaler(m cqrs.CommandEventMarshaler) eventGroupProcessorOption {
	return eventGroupProcessorOptionFunc(func(c *cqrs.EventGroupProcessorConfig) { c.Marshaler = m })
}

func WithEventGroupProcessorGenerateSubscribeTopic(fn cqrs.EventGroupProcessorGenerateSubscribeTopicFn) eventGroupProcessorOption {
	return eventGroupProcessorOptionFunc(func(c *cqrs.EventGroupProcessorConfig) { c.GenerateSubscribeTopic = fn })
}

func defaultEventGroupProcessorGenerateSubscribeTopic(params cqrs.EventGroupProcessorGenerateSubscribeTopicParams) (string, error) {
	return "eventGroups-" + kebabcase.Kebabcase(params.EventGroupName), nil
}

func defaultEventGroupProcessorConfig() cqrs.EventGroupProcessorConfig {
	return cqrs.EventGroupProcessorConfig{
		GenerateSubscribeTopic: defaultEventGroupProcessorGenerateSubscribeTopic,
		Marshaler:              DefaultMarshaler(),
		Logger:                 watermill.NewSlogLogger(nil),
	}
}
