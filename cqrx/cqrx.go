package cqrx

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/cyg-pd/go-watermillx/driver"
	"github.com/cyg-pd/go-watermillx/internal/utils"
)

type CQRS struct {
	router *message.Router
	driver driver.Driver

	publisherDecorator  []message.PublisherDecorator
	subscriberDecorator []message.SubscriberDecorator
}

func (c *CQRS) AddPublisherDecorator(decorator message.PublisherDecorator) {
	c.publisherDecorator = append(c.publisherDecorator, decorator)
}

func (c *CQRS) AddSubscriberDecorator(decorator message.SubscriberDecorator) {
	c.subscriberDecorator = append(c.subscriberDecorator, decorator)
}

func (c *CQRS) applyPublisherDecorator(pub message.Publisher) (message.Publisher, error) {
	var err error
	for _, decorator := range c.publisherDecorator {
		if pub, err = decorator(pub); err != nil {
			return nil, err
		}
	}
	return pub, nil
}

func (c *CQRS) applySubscriberDecorator(sub message.Subscriber) (message.Subscriber, error) {
	var err error
	for _, decorator := range c.subscriberDecorator {
		if sub, err = decorator(sub); err != nil {
			return nil, err
		}
	}
	return sub, nil
}

func (c *CQRS) buildPublisher() (message.Publisher, error) {
	pub, err := c.driver.Publisher()
	if err != nil {
		return nil, fmt.Errorf("watermillx/cqrx: failed to initialize publisher driver: %w", err)
	}

	if pub, err = c.applyPublisherDecorator(pub); err != nil {
		return nil, err
	}

	return pub, nil
}

func (c *CQRS) initializeTopic(sub message.Subscriber, topic string) error {
	if sub, ok := sub.(message.SubscribeInitializer); ok {
		if err := sub.SubscribeInitialize(topic); err != nil {
			return fmt.Errorf("watermillx/cqrx: failed to initialize topic: %w", err)
		}
	}
	return nil
}

func (c *CQRS) buildSubscriber(topic string, opts ...driver.SubscriberOption) (message.Subscriber, error) {
	sub, err := c.driver.Subscriber(opts...)
	if err != nil {
		return nil, fmt.Errorf("watermillx/cqrx: failed to initialize subscriber driver: %w", err)
	}

	if err = c.initializeTopic(sub, topic); err != nil {
		return nil, err
	}

	if sub, err = c.applySubscriberDecorator(sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (c *CQRS) MustCommandBus(opts ...commandBusOption) *cqrs.CommandBus {
	return utils.Must(c.CommandBus(opts...))
}

func (c *CQRS) CommandBus(opts ...commandBusOption) (*cqrs.CommandBus, error) {
	pub, err := c.buildPublisher()
	if err != nil {
		return nil, err
	}

	conf := defaultCommandBusConfig()
	for _, opt := range opts {
		opt.apply(&conf)
	}

	v, err := cqrs.NewCommandBusWithConfig(pub, conf)
	if err != nil {
		return nil, fmt.Errorf("watermillx/cqrx: failed to create cqrs.CommandBus: %w", err)
	}
	return v, nil
}

func (c *CQRS) MustCommandProcessor(opts ...commandProcessorOption) *cqrs.CommandProcessor {
	return utils.Must(c.CommandProcessor(opts...))
}

func (c *CQRS) CommandProcessor(opts ...commandProcessorOption) (*cqrs.CommandProcessor, error) {
	conf := defaultCommandProcessorConfig()
	for _, opt := range opts {
		opt.apply(&conf)
	}

	if conf.GenerateSubscribeTopic == nil {
		conf.GenerateSubscribeTopic = defaultCommandProcessorGenerateSubscribeTopic
	}

	conf.SubscriberConstructor = func(params cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
		topic, err := conf.GenerateSubscribeTopic(cqrs.CommandProcessorGenerateSubscribeTopicParams{
			CommandName:    params.CommandName,
			CommandHandler: params.Handler,
		})
		if err != nil {
			return nil, err
		}

		if h, ok := params.Handler.(HandlerConsumerGroup); ok {
			if g := h.ConsumerGroup(); g != "" {
				return c.buildSubscriber(topic, driver.WithSubscriberConsumerGroup(g))
			}
		}
		return c.buildSubscriber(topic, driver.WithSubscriberConsumerGroup(topic))
	}

	v, err := cqrs.NewCommandProcessorWithConfig(c.router, conf)
	if err != nil {
		return nil, fmt.Errorf("watermillx/cqrx: failed to create cqrs.CommandProcessor: %w", err)
	}
	return v, nil
}

func (c *CQRS) MustEventBus(opts ...eventBusOption) *cqrs.EventBus {
	return utils.Must(c.EventBus(opts...))
}

func (c *CQRS) EventBus(opts ...eventBusOption) (*cqrs.EventBus, error) {
	pub, err := c.buildPublisher()
	if err != nil {
		return nil, err
	}

	conf := defaultEventBusConfig()
	for _, opt := range opts {
		opt.apply(&conf)
	}

	v, err := cqrs.NewEventBusWithConfig(pub, conf)
	if err != nil {
		return nil, fmt.Errorf("watermillx/cqrx: failed to create cqrs.EventBus: %w", err)
	}
	return v, nil
}

func (c *CQRS) MustEventProcessor(opts ...eventProcessorOption) *cqrs.EventProcessor {
	return utils.Must(c.EventProcessor(opts...))
}

func (c *CQRS) EventProcessor(opts ...eventProcessorOption) (*cqrs.EventProcessor, error) {
	conf := defaultEventProcessorConfig()
	for _, opt := range opts {
		opt.apply(&conf)
	}

	if conf.GenerateSubscribeTopic == nil {
		conf.GenerateSubscribeTopic = defaultEventProcessorGenerateSubscribeTopic
	}

	conf.SubscriberConstructor = func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
		topic, err := conf.GenerateSubscribeTopic(cqrs.EventProcessorGenerateSubscribeTopicParams{
			EventName:    params.EventName,
			EventHandler: params.EventHandler,
		})
		if err != nil {
			return nil, err
		}

		if h, ok := params.EventHandler.(HandlerConsumerGroup); ok {
			if g := h.ConsumerGroup(); g != "" {
				return c.buildSubscriber(topic, driver.WithSubscriberConsumerGroup(g))
			}
		}
		return c.buildSubscriber(topic, driver.WithSubscriberConsumerGroup(topic))

	}

	v, err := cqrs.NewEventProcessorWithConfig(c.router, conf)
	if err != nil {
		return nil, fmt.Errorf("watermillx/cqrx: failed to create cqrs.EventProcessor: %w", err)
	}
	return v, nil
}

func (c *CQRS) MustEventGroupProcessor(opts ...eventGroupProcessorOption) *cqrs.EventGroupProcessor {
	return utils.Must(c.EventGroupProcessor(opts...))
}

func (c *CQRS) EventGroupProcessor(opts ...eventGroupProcessorOption) (*cqrs.EventGroupProcessor, error) {
	conf := defaultEventGroupProcessorConfig()
	for _, opt := range opts {
		opt.apply(&conf)
	}

	if conf.GenerateSubscribeTopic == nil {
		conf.GenerateSubscribeTopic = defaultEventGroupProcessorGenerateSubscribeTopic
	}

	conf.SubscriberConstructor = func(params cqrs.EventGroupProcessorSubscriberConstructorParams) (message.Subscriber, error) {
		topic, err := conf.GenerateSubscribeTopic(cqrs.EventGroupProcessorGenerateSubscribeTopicParams(params))
		if err != nil {
			return nil, err
		}

		return c.buildSubscriber(topic, driver.WithSubscriberConsumerGroup(params.EventGroupName))

	}

	v, err := cqrs.NewEventGroupProcessorWithConfig(c.router, conf)
	if err != nil {
		return nil, fmt.Errorf("watermillx/cqrx: failed to create cqrs.EventProcessor: %w", err)
	}
	return v, nil
}

func New(router *message.Router, driver driver.Driver) *CQRS {
	return &CQRS{
		router: router,
		driver: driver,
	}
}
