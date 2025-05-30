package driver

import (
	"fmt"
	"sync"
)

var store sync.Map

func Register(name string, driver DriverFunc) {
	store.Store(name, driver)
}

func New(name string, config any) (Driver, error) {
	if d, ok := store.Load(name); ok {
		return d.(DriverFunc)(config)
	}

	return nil, fmt.Errorf("watermillx/driver: %s driver not found", name)
}
