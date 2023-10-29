package dgate

import (
	"dgate/types"
	"fmt"
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
	c.AddHandler(types.EventNameReady, func(e *types.ReadyEventData) {
		fmt.Println(e)
	})
	c.AddHandler(types.EventNameMessageCreate, func(e *types.MessageEventData) {
		fmt.Println(e)
	})
	handleErr(c.Connect())
}
