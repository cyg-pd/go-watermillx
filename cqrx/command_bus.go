package cqrx

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/cyg-pd/go-kebabcase"
)

type commandBusOption interface{ apply(*cqrs.CommandBusConfig) }
type commandBusOptionFunc func(*cqrs.CommandBusConfig)

func (fn commandBusOptionFunc) apply(cfg *cqrs.CommandBusConfig) { fn(cfg) }

func WithCommandBusOnSend(fn func(params cqrs.CommandBusOnSendParams) error) commandBusOption {
	return commandBusOptionFunc(func(c *cqrs.CommandBusConfig) { c.OnSend = fn })
}

func WithCommandBusMarshaler(m cqrs.CommandEventMarshaler) commandBusOption {
	return commandBusOptionFunc(func(c *cqrs.CommandBusConfig) { c.Marshaler = m })
}

func WithCommandBusGeneratePublishTopic(fn cqrs.CommandBusGeneratePublishTopicFn) commandBusOption {
	return commandBusOptionFunc(func(c *cqrs.CommandBusConfig) { c.GeneratePublishTopic = fn })
}

func defaultCommandBusGeneratePublishTopic(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
	return "commands-" + kebabcase.Kebabcase(params.CommandName), nil
}

func defaultCommandBusConfig() cqrs.CommandBusConfig {
	return cqrs.CommandBusConfig{
		GeneratePublishTopic: defaultCommandBusGeneratePublishTopic,
		Marshaler:            DefaultMarshaler(),
		Logger:               watermill.NewSlogLogger(nil),
	}
}
