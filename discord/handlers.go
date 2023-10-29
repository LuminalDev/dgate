package discord

import (
	"dgate/types"
	"fmt"
	"sync"
)

type Handlers struct {
	OnReady         []func(data *types.ReadyEventData)
	OnMessageCreate []func(data *types.MessageEventData)
	mutex           sync.Mutex
}

func (handlers *Handlers) Add(event string, function any) error {
	handlers.mutex.Lock()
	defer handlers.mutex.Unlock()

	switch event {
	case types.EventNameReady:
		if function, ok := function.(func(data *types.ReadyEventData)); ok {
			handlers.OnReady = append(handlers.OnReady, function)
		} else {
			return fmt.Errorf("wrong function signature")
		}
	case types.EventNameMessageCreate:
		if function, ok := function.(func(data *types.MessageEventData)); ok {
			handlers.OnMessageCreate = append(handlers.OnMessageCreate, function)
		} else {
			return fmt.Errorf("wrong function signature")
		}
	}
	return nil
}
