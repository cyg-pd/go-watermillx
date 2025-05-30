package driver

type SubscriberConfig struct {
	ConsumerGroup string
}

func (c *SubscriberConfig) Apply(opts []SubscriberOption) {
	for _, opt := range opts {
		opt.apply(c)
	}
}

type SubscriberOption interface{ apply(*SubscriberConfig) }
type SubscriberOptionFunc func(*SubscriberConfig)

func (fn SubscriberOptionFunc) apply(cfg *SubscriberConfig) { fn(cfg) }

func WithSubscriberConsumerGroup(g string) SubscriberOption {
	return SubscriberOptionFunc(func(c *SubscriberConfig) { c.ConsumerGroup = g })
}
