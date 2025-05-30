package manager

import (
	"fmt"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/cyg-pd/go-watermillx/cqrx"
	"github.com/cyg-pd/go-watermillx/driver"
	"github.com/cyg-pd/go-watermillx/internal/utils"
)

type Manager struct {
	store  sync.Map
	router *message.Router
}

func (m *Manager) Add(name string, driver driver.Driver) {
	m.store.Store(name, driver)
}

func (m *Manager) MustUse(name string) driver.Driver { return utils.Must(m.Use(name)) }
func (m *Manager) Use(name string) (driver.Driver, error) {
	dr, ok := m.store.Load(name)
	if !ok {
		return nil, fmt.Errorf("watermillx/manager: %s not exist", name)
	}
	return dr.(driver.Driver), nil
}

func (m *Manager) MustUseCQRS(name string) *cqrx.CQRS { return utils.Must(m.UseCQRS(name)) }
func (m *Manager) UseCQRS(name string) (*cqrx.CQRS, error) {
	dr, err := m.Use(name)
	if err != nil {
		return nil, err
	}
	return cqrx.New(
		m.router,
		dr,
	), nil
}

func New(router *message.Router) *Manager {
	return &Manager{
		router: router,
	}
}
