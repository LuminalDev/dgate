package discord

import (
	"errors"
	"sync"

	"github.com/luminaldev/dgate/types"
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

	failed := false

	switch event {
	case types.ReadyEventHandler:
		if function, ok := function.(func(data *types.ReadyEventData)); ok {
			handlers.OnReady = append(handlers.OnReady, function)
		} else {
			failed = true
		}
	case types.MessageCreateEventHandler:
		if function, ok := function.(func(data *types.MessageEventData)); ok {
			handlers.OnMessageCreate = append(handlers.OnMessageCreate, function)
		} else {
			failed = true
		}
	case types.MessageUpdateEventHandler:
		if function, ok := function.(func(data *types.MessageEventData)); ok {
			handlers.OnMessageUpdate = append(handlers.OnMessageUpdate, function)
		} else {
			failed = true
		}
	case types.InvalidatedEventHandler:
		if function, ok := function.(func()); ok {
			handlers.OnInvalidated = append(handlers.OnInvalidated, function)
		} else {
			failed = true
		}
	case types.ReconnectEventHandler:
		if function, ok := function.(func()); ok {
			handlers.OnReconnect = append(handlers.OnReconnect, function)
		} else {
			failed = true
		}
	default:
		return errors.New("failed to match event to gateway event")
	}

	if failed {
		return errors.New("function signature was not correct for the specified event")
	}

	return nil
}
