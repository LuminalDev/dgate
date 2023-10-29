package dgate

import (
	"fmt"
	"os"
	"testing"

	"github.com/luminaldev/dgate/types"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Test_Main(t *testing.T) {
	token, err := os.ReadFile("test_token")
	handleErr(err)

	client := NewClient(string(token))

	client.AddHandler(types.GatewayEventReady, func(e *types.ReadyEventData) {
		fmt.Printf("%#v", e)
	})
	client.AddHandler(types.GatewayEventMessageCreate, func(e *types.MessageEventData) {
		fmt.Printf("%#v", e)
	})
	client.AddHandler(types.GatewayEventMessageUpdate, func(e *types.MessageEventData) {
		fmt.Printf("%#v", e)
	})

	handleErr(client.Connect())
}
