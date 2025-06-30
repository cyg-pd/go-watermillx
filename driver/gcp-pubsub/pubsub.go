package pubsub

import (
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	watermillgooglecloud "github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/cyg-pd/go-watermillx/driver"
	"github.com/go-viper/mapstructure/v2"
)

func init() {
	driver.Register(
		"gcp-pubsub",
		func(c any) (driver.Driver, error) { return New(c) },
	)
}

type GoogleCloudPubSub struct {
	conf *GoogleCloudPubSubConfig
}

func (g *GoogleCloudPubSub) Publisher() (message.Publisher, error) {
	conf := watermillgooglecloud.PublisherConfig{
		ProjectID:                                     g.conf.ProjectID,
		DoNotCheckTopicExistence:                      g.conf.DoNotCheckTopicExistence,
		DoNotCreateTopicIfMissing:                     g.conf.DoNotCreateTopicIfMissing,
		EnableMessageOrdering:                         g.conf.Publisher.EnableMessageOrdering,
		EnableMessageOrderingAutoResumePublishOnError: g.conf.Publisher.EnableMessageOrderingAutoResumePublishOnError,
	}

	if g.conf.TopicProjectID != "" {
		conf.ProjectID = g.conf.TopicProjectID
	}

	publisher, err := watermillgooglecloud.NewPublisher(conf, watermill.NewSlogLogger(nil))
	if err != nil {
		return nil, fmt.Errorf("watermillx/driver/gcp-pubsub: create new publisher error: %w", err)
	}

	return publisher, nil
}

func (g *GoogleCloudPubSub) Subscriber(opts ...driver.SubscriberOption) (message.Subscriber, error) {
	var c driver.SubscriberConfig
	c.Apply(opts)

	conf := watermillgooglecloud.SubscriberConfig{
		GenerateSubscriptionName: func(topic string) string {
			if c.ConsumerGroup != "" {
				return c.ConsumerGroup
			}
			return topic
		},
		ProjectID:                        g.conf.ProjectID,
		TopicProjectID:                   g.conf.TopicProjectID,
		DoNotCreateTopicIfMissing:        g.conf.DoNotCreateTopicIfMissing,
		DoNotCreateSubscriptionIfMissing: g.conf.DoNotCreateSubscriptionIfMissing,
		ReceiveSettings:                  g.conf.Subscriber.ReceiveSettings,
		SubscriptionConfig:               g.conf.Subscriber.SubscriptionConfig,
	}

	subscriber, err := watermillgooglecloud.NewSubscriber(conf, watermill.NewSlogLogger(nil))
	if err != nil {
		return nil, fmt.Errorf("watermillx/driver/gcp-pubsub: create new subscriber error: %w", err)
	}

	return subscriber, nil
}

func New(config any) (*GoogleCloudPubSub, error) {
	c := NewGoogleCloudPubSubConfig()

	switch v := config.(type) {
	case *GoogleCloudPubSubConfig:
		c = v
	case GoogleCloudPubSubConfig:
		c = &v
	case map[string]any:
		if err := mapstructure.Decode(v, c); err != nil {
			return nil, fmt.Errorf("watermillx/driver/gcp-pubsub: parse config error: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(v), c); err != nil {
			return nil, fmt.Errorf("watermillx/driver/gcp-pubsub: parse config error: %w", err)
		}
	case []byte:
		if err := json.Unmarshal(v, c); err != nil {
			return nil, fmt.Errorf("watermillx/driver/gcp-pubsub: parse config error: %w", err)
		}
	default:
		return nil, fmt.Errorf("watermillx/driver/gcp-pubsub: config only support GoogleCloudPubSubConfig, map[string]any, string, []byte")
	}

	return &GoogleCloudPubSub{conf: c}, nil
}

var _ driver.Driver = (*GoogleCloudPubSub)(nil)
