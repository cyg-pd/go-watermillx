package pubsub

import "cloud.google.com/go/pubsub"

type GoogleCloudPubSubConfig struct {
	ProjectID                        string
	TopicProjectID                   string
	DoNotCheckTopicExistence         bool
	DoNotCreateTopicIfMissing        bool
	DoNotCreateSubscriptionIfMissing bool
	Subscriber                       struct {
		ReceiveSettings    pubsub.ReceiveSettings
		SubscriptionConfig pubsub.SubscriptionConfig
	}
	Publisher struct {
		EnableMessageOrdering                         bool
		EnableMessageOrderingAutoResumePublishOnError bool
	}
}

func NewGoogleCloudPubSubConfig() *GoogleCloudPubSubConfig {
	return &GoogleCloudPubSubConfig{}
}
