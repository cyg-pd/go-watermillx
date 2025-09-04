package pubsub

import (
	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)

type GoogleCloudPubSubConfig struct {
	ProjectID                        string
	TopicProjectID                   string
	DoNotCheckTopicExistence         bool
	DoNotCreateTopicIfMissing        bool
	DoNotCreateSubscriptionIfMissing bool
	Subscriber                       struct {
		ReceiveSettings    pubsub.ReceiveSettings
		SubscriptionConfig pubsubpb.Subscription
	}
	Publisher struct {
		EnableMessageOrdering                         bool
		EnableMessageOrderingAutoResumePublishOnError bool
	}
}

func NewGoogleCloudPubSubConfig() *GoogleCloudPubSubConfig {
	return &GoogleCloudPubSubConfig{}
}
