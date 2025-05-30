package cqrx

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/cyg-pd/go-kebabcase"
)

type commandProcessorOption interface {
	apply(*cqrs.CommandProcessorConfig)
}
type commandProcessorOptionFunc func(*cqrs.CommandProcessorConfig)

func (fn commandProcessorOptionFunc) apply(cfg *cqrs.CommandProcessorConfig) { fn(cfg) }

func WithCommandProcessorOnHandle(fn func(params cqrs.CommandProcessorOnHandleParams) error) commandProcessorOption {
	return commandProcessorOptionFunc(func(c *cqrs.CommandProcessorConfig) { c.OnHandle = fn })
}

func WithCommandProcessorMarshaler(m cqrs.CommandEventMarshaler) commandProcessorOption {
	return commandProcessorOptionFunc(func(c *cqrs.CommandProcessorConfig) { c.Marshaler = m })
}

func WithCommandProcessorGenerateSubscribeTopic(fn cqrs.CommandProcessorGenerateSubscribeTopicFn) commandProcessorOption {
	return commandProcessorOptionFunc(func(c *cqrs.CommandProcessorConfig) { c.GenerateSubscribeTopic = fn })
}

func defaultCommandProcessorGenerateSubscribeTopic(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
	return "commands-" + kebabcase.Kebabcase(params.CommandName), nil
}

func defaultCommandProcessorConfig() cqrs.CommandProcessorConfig {
	return cqrs.CommandProcessorConfig{
		GenerateSubscribeTopic: defaultCommandProcessorGenerateSubscribeTopic,
		Marshaler:              DefaultMarshaler(),
		Logger:                 watermill.NewSlogLogger(nil),
	}
}
