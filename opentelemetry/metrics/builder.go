package metrics

import "github.com/ThreeDotsLabs/watermill/message"

type Builder interface {
	RouterMiddleware(h message.HandlerFunc) message.HandlerFunc
	DecoratePublisher(pub message.Publisher) (message.Publisher, error)
	DecorateSubscriber(sub message.Subscriber) (message.Subscriber, error)
}
