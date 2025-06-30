package redis

import (
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/cyg-pd/go-watermillx/driver"
	"github.com/go-viper/mapstructure/v2"
	"github.com/redis/go-redis/v9"
)

func init() {
	driver.Register(
		"redis",
		func(c any) (driver.Driver, error) { return New(c) },
	)
}

type Redis struct {
	client redis.UniversalClient
	conf   *RedisConfig
}

func (r *Redis) Publisher() (message.Publisher, error) {
	return redisstream.NewPublisher(
		redisstream.PublisherConfig{
			Client: r.client,
		},
		nil,
	)
}

func (r *Redis) Subscriber(opts ...driver.SubscriberOption) (message.Subscriber, error) {
	var c driver.SubscriberConfig
	c.Apply(opts)

	return redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			ConsumerGroup:          c.ConsumerGroup,
			Client:                 r.client,
			NackResendSleep:        r.conf.Subscriber.NackResendSleep,
			BlockTime:              r.conf.Subscriber.BlockTime,
			ClaimInterval:          r.conf.Subscriber.ClaimInterval,
			ClaimBatchSize:         r.conf.Subscriber.ClaimBatchSize,
			MaxIdleTime:            r.conf.Subscriber.MaxIdleTime,
			CheckConsumersInterval: r.conf.Subscriber.CheckConsumersInterval,
			ConsumerTimeout:        r.conf.Subscriber.ConsumerTimeout,
		},
		nil,
	)
}

func New(config any) (*Redis, error) {
	c := NewRedisConfig()

	switch v := config.(type) {
	case *RedisConfig:
		c = v
	case RedisConfig:
		c = &v
	case map[string]any:
		if err := mapstructure.Decode(v, c); err != nil {
			return nil, fmt.Errorf("watermillx/driver/redis: parse config error: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(v), c); err != nil {
			return nil, fmt.Errorf("watermillx/driver/redis: parse config error: %w", err)
		}
	case []byte:
		if err := json.Unmarshal(v, c); err != nil {
			return nil, fmt.Errorf("watermillx/driver/redis: parse config error: %w", err)
		}
	default:
		return nil, fmt.Errorf("watermillx/driver/redis: config only support RedisConfig, map[string]any, string, []byte")
	}

	url, err := redis.ParseURL(c.DSN)
	if err != nil {
		return nil, fmt.Errorf("watermillx/driver/redis: parse DSN error: %w", err)
	}

	return &Redis{conf: c, client: redis.NewClient(url)}, nil
}

var _ driver.Driver = (*Redis)(nil)
