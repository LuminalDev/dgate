package dgate

import (
	"dgate/discord"
	"fmt"
	"github.com/switchupcb/dasgo/dasgo"
	"os"
	"testing"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Test_Main(t *testing.T) {
	token, err := os.ReadFile("test_token")
	handleErr(err)
	c := NewClient(string(token))
	c.AddHandler(dasgo.FlagGatewayEventNameReady, func(e *discord.ReadyData) {
		fmt.Println(e)
	})
	c.AddHandler(dasgo.FlagGatewayEventNameMessageCreate, func(e *discord.MessageData) {
		fmt.Println(e)
	})
	handleErr(c.Connect())
}
