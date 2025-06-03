package kafka

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/IBM/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/cyg-pd/go-watermillx/driver"
	"github.com/go-viper/mapstructure/v2"
)

func init() {
	driver.Register(
		"kafka",
		func(c any) (driver.Driver, error) { return New(c) },
	)
}

type Kafka struct {
	conf *KafkaConfig
}

func (k *Kafka) Publisher() (message.Publisher, error) {
	conf := kafka.PublisherConfig{
		Brokers:               k.conf.Brokers,
		OverwriteSaramaConfig: k.conf.SaramaConfig,
		Marshaler:             kafka.DefaultMarshaler{},
	}

	pub, err := kafka.NewPublisher(conf, watermill.NewSlogLogger(nil))
	if err != nil {
		return nil, fmt.Errorf("watermillx/driver/kafka: create new publisher error: %w", err)
	}

	return pub, nil
}

func (k *Kafka) Subscriber(opts ...driver.SubscriberOption) (message.Subscriber, error) {
	conf := kafka.SubscriberConfig{
		Brokers:                k.conf.Brokers,
		OverwriteSaramaConfig:  k.conf.SaramaConfig,
		ConsumerGroup:          k.conf.Subscriber.ConsumerGroup,
		InitializeTopicDetails: k.conf.InitializeTopicDetails,
		Unmarshaler:            kafka.DefaultMarshaler{},
	}

	var c driver.SubscriberConfig
	c.Apply(opts)
	if c.ConsumerGroup != "" {
		conf.ConsumerGroup = c.ConsumerGroup
	}

	if conf.InitializeTopicDetails == nil {
		conf.InitializeTopicDetails = &sarama.TopicDetail{
			NumPartitions:     -1,
			ReplicationFactor: -1,
		}
	}

	if conf.InitializeTopicDetails.NumPartitions == 0 {
		conf.InitializeTopicDetails.NumPartitions = -1
	}

	if conf.InitializeTopicDetails.ReplicationFactor == 0 {
		conf.InitializeTopicDetails.ReplicationFactor = -1
	}

	pub, err := kafka.NewSubscriber(conf, watermill.NewSlogLogger(nil))
	if err != nil {
		return nil, fmt.Errorf("watermillx/driver/kafka: create new subscriber error: %w", err)
	}

	return pub, nil
}

func New(config any) (*Kafka, error) {
	c := NewKafkaConfig()

	switch v := config.(type) {
	case *KafkaConfig:
		c = v
	case KafkaConfig:
		c = &v
	case map[string]any:
		if err := mapstructure.Decode(v, c); err != nil {
			return nil, fmt.Errorf("watermillx/driver/kafka: parse config error: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(v), c); err != nil {
			return nil, fmt.Errorf("watermillx/driver/kafka: parse config error: %w", err)
		}
	case []byte:
		if err := json.Unmarshal(v, c); err != nil {
			return nil, fmt.Errorf("watermillx/driver/kafka: parse config error: %w", err)
		}
	default:
		return nil, fmt.Errorf("watermillx/driver/kafka: config only support KafkaConfig, map[string]any, string, []byte")
	}

	if c.SaramaConfig == nil {
		c.SaramaConfig = DefaultSaramaConfig()
	}

	if v, err := os.Hostname(); err == nil {
		if c.SaramaConfig.ClientID == "watermill" || c.SaramaConfig.ClientID == "" {
			c.SaramaConfig.ClientID = v
		}
	}

	return &Kafka{conf: c}, nil
}

var _ driver.Driver = (*Kafka)(nil)
