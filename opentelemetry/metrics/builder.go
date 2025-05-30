package metrics

import "github.com/ThreeDotsLabs/watermill/message"

type Builder interface {
	AttachRouter(router *message.Router)
	DecoratePublisher(pub message.Publisher) (message.Publisher, error)
	DecorateSubscriber(sub message.Subscriber) (message.Subscriber, error)
}
