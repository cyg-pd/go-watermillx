package kafka

import (
	"github.com/IBM/sarama"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
)

type KafkaConfig struct {
	// Kafka brokers list.
	Brokers []string

	// SaramaConfig holds additional sarama settings.
	SaramaConfig *sarama.Config

	Subscriber struct {
		// Kafka consumer group.
		// When empty, all messages from all partitions will be returned.
		ConsumerGroup string
	}

	InitializeTopicDetails *sarama.TopicDetail
}

func NewKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		SaramaConfig: DefaultSaramaConfig(),
	}
}

func DefaultSaramaConfig() *sarama.Config {
	c := kafka.DefaultSaramaSyncPublisherConfig()
	c.Producer.RequiredAcks = sarama.WaitForAll
	c.Producer.Return.Successes = true
	c.Consumer.Return.Errors = true
	return c
}
