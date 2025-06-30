package redis

import "time"

type RedisConfig struct {
	DSN string `mapstructure:"dsn" json:"dsn"`

	// Subscriber redisstream.SubscriberConfig
	Subscriber struct {
		NackResendSleep        time.Duration
		BlockTime              time.Duration
		ClaimInterval          time.Duration
		ClaimBatchSize         int64
		MaxIdleTime            time.Duration
		CheckConsumersInterval time.Duration
		ConsumerTimeout        time.Duration
	}
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{}
}
