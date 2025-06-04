package messagex

import (
	"context"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
)

// MessageTransformSubscriberDecorator creates a subscriber decorator that calls transform
// on each message that passes through the subscriber.
func MessageTransformSubscriberDecorator(transform func(topic string, msg *message.Message)) message.SubscriberDecorator {
	if transform == nil {
		panic("transform function is nil")
	}
	return func(sub message.Subscriber) (message.Subscriber, error) {
		return &messageTransformSubscriberDecorator{
			sub:       sub,
			transform: transform,
		}, nil
	}
}

type messageTransformSubscriberDecorator struct {
	sub message.Subscriber

	transform   func(topic string, msg *message.Message)
	subscribeWg sync.WaitGroup
}

func (t *messageTransformSubscriberDecorator) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	in, err := t.sub.Subscribe(ctx, topic)
	if err != nil {
		return nil, err
	}

	out := make(chan *message.Message)
	t.subscribeWg.Add(1)
	go func() {
		for msg := range in {
			t.transform(topic, msg)
			out <- msg
		}
		close(out)
		t.subscribeWg.Done()
	}()

	return out, nil
}

func (t *messageTransformSubscriberDecorator) Close() error {
	err := t.sub.Close()

	t.subscribeWg.Wait()
	return err
}
