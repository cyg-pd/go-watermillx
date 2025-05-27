package driver

import (
	"github.com/ThreeDotsLabs/watermill/message"
)

type DriverFunc func(config any) (Driver, error)
type Driver interface {
	Publisher() (message.Publisher, error)
	Subscriber(topic string, opts ...SubscriberOption) (message.Subscriber, error)
}
