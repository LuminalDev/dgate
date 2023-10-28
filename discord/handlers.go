package discord

import (
	"fmt"
	"github.com/switchupcb/dasgo/dasgo"
	"sync"
)

type Handlers struct {
	OnReady         []func(*ReadyData)
	OnMessageCreate []func(data *MessageData)
	mutex           sync.Mutex
}

func (handlers *Handlers) Add(event string, function any) error {
	handlers.mutex.Lock()
	defer handlers.mutex.Unlock()

	switch event {
	case dasgo.FlagGatewayEventNameReady:
		if function, ok := function.(func(*ReadyData)); ok {
			handlers.OnReady = append(handlers.OnReady, function)
		} else {
			return fmt.Errorf("wrong function signature")
		}
	case dasgo.FlagGatewayEventNameMessageCreate:
		if function, ok := function.(func(data *MessageData)); ok {
			handlers.OnMessageCreate = append(handlers.OnMessageCreate, function)
		} else {
			return fmt.Errorf("wrong function signature")
		}
	}
	return nil
}
