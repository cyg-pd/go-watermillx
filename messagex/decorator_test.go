package messagex_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/pubsub/tests"
	"github.com/cyg-pd/go-watermillx/messagex"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/subscriber"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

var noop = func(string, *message.Message) {}
var errClosing = errors.New("mock error on close")

type mockSubscriber struct {
	ch chan *message.Message
}

func (m mockSubscriber) Subscribe(context.Context, string) (<-chan *message.Message, error) {
	return m.ch, nil
}

func (m mockSubscriber) Close() error {
	close(m.ch)
	return nil
}

func TestMessageTransformSubscriberDecorator_transparent(t *testing.T) {
	sub := mockSubscriber{make(chan *message.Message)}
	decorated, err := messagex.MessageTransformSubscriberDecorator(noop)(sub)
	require.NoError(t, err)

	messages, err := decorated.Subscribe(context.Background(), "topic")
	require.NoError(t, err)

	richMessage := message.NewMessage("uuid", []byte("serious payloads"))
	richMessage.Metadata.Set("k1", "v1")
	richMessage.Metadata.Set("k2", "v2")

	go func() {
		sub.ch <- richMessage
	}()

	received, all := subscriber.BulkRead(messages, 1, time.Second)
	require.True(t, all)

	assert.True(t, received[0].Equals(richMessage), "expected the message to pass unchanged through decorator")
}

type closingSubscriber struct {
	closed bool
}

func (closingSubscriber) Subscribe(context.Context, string) (<-chan *message.Message, error) {
	return nil, nil
}

func (c *closingSubscriber) Close() error {
	c.closed = true
	return errClosing
}

func TestMessageTransformSubscriberDecorator_Close(t *testing.T) {
	cs := &closingSubscriber{}

	decoratedSub, err := messagex.MessageTransformSubscriberDecorator(noop)(cs)
	require.NoError(t, err)

	// given
	require.False(t, cs.closed)

	// when
	decoratedCloseErr := decoratedSub.Close()

	// then
	assert.True(
		t,
		cs.closed,
		"expected the Close() call to propagate to decorated subscriber",
	)
	assert.Equal(
		t,
		errClosing,
		decoratedCloseErr,
		"expected the decorator to propagate the closing error from underlying subscriber",
	)
}

func TestMessageTransformSubscriberDecorator_Subscribe(t *testing.T) {
	numMessages := 1000
	pubSub := gochannel.NewGoChannel(gochannel.Config{}, watermill.NewStdLogger(true, true))

	onMessage := func(topic string, msg *message.Message) {
		msg.Metadata.Set("key", topic)
	}
	decorator := messagex.MessageTransformSubscriberDecorator(onMessage)

	decoratedSub, err := decorator(pubSub)
	require.NoError(t, err)

	messages, err := decoratedSub.Subscribe(context.Background(), "topic")
	require.NoError(t, err)

	sent := message.Messages{}

	go func() {
		for i := 0; i < numMessages; i++ {
			msg := message.NewMessage(strconv.Itoa(i), []byte{})
			sent = append(sent, msg)

			err = pubSub.Publish("topic", msg)
			require.NoError(t, err)
		}
	}()

	received, all := subscriber.BulkRead(messages, numMessages, time.Second)
	require.True(t, all)
	tests.AssertAllMessagesReceived(t, sent, received)

	for _, msg := range received {
		assert.Equal(
			t,
			"topic",
			msg.Metadata.Get("key"),
			"expected onMessage callback to have set metadata",
		)
	}
}

func TestMessageTransformer_nil_panics(t *testing.T) {
	require.Panics(
		t,
		func() {
			_ = messagex.MessageTransformSubscriberDecorator(nil)
		},
		"expected to panic if transform is nil",
	)
}
