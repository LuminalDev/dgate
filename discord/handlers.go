package discord

import (
	"fmt"
	"github.com/luminaldev/dgate/types"
	"sync"
)

type Handlers struct {
	OnReady         []func(data *types.ReadyEventData)
	OnMessageCreate []func(data *types.MessageEventData)
	OnMessageUpdate []func(data *types.MessageEventData)
	OnReconnect     []func()
	mutex           sync.Mutex
	OnInvalidated   []func()
}

func (handlers *Handlers) Add(event string, function any) error {
	handlers.mutex.Lock()
	defer handlers.mutex.Unlock()

	switch event {
	case types.ReadyEventHandler:
		if function, ok := function.(func(data *types.ReadyEventData)); ok {
			handlers.OnReady = append(handlers.OnReady, function)
		} else {
			return fmt.Errorf("wrong function signature")
		}
	case types.MessageCreateEventHandler:
		if function, ok := function.(func(data *types.MessageEventData)); ok {
			handlers.OnMessageCreate = append(handlers.OnMessageCreate, function)
		} else {
			return fmt.Errorf("wrong function signature")
		}
	case types.MessageUpdateEventHandler:
		if function, ok := function.(func(data *types.MessageEventData)); ok {
			handlers.OnMessageUpdate = append(handlers.OnMessageUpdate, function)
		} else {
			return fmt.Errorf("wrong function signature")
		}
	case types.InvalidatedEventHandler:
		if function, ok := function.(func()); ok {
			handlers.OnInvalidated = append(handlers.OnInvalidated, function)
		} else {
			return fmt.Errorf("wrong function signature")
		}
	case types.ReconnectEventHandler:
		if function, ok := function.(func()); ok {
			handlers.OnReconnect = append(handlers.OnReconnect, function)
		} else {
			return fmt.Errorf("wrong function signature")
		}
	}
	return nil
}
