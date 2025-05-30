package router

import (
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type option interface{ apply(*message.RouterConfig) }
type optionFunc func(*message.RouterConfig)

func (fn optionFunc) apply(cfg *message.RouterConfig) { fn(cfg) }

func WithCloseTimeout(t time.Duration) option {
	return optionFunc(func(cfg *message.RouterConfig) { cfg.CloseTimeout = t })
}

func New(opts ...option) (*message.Router, error) {
	var c message.RouterConfig
	for _, opt := range opts {
		opt.apply(&c)
	}

	return message.NewRouter(c, watermill.NewSlogLogger(nil))
}
