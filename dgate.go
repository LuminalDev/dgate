package dgate

import "github.com/luminaldev/dgate/discord"

type Client struct {
	Selfbot *discord.Selfbot
	Gateway *discord.Gateway
}

func NewClient(token string) *Client {
	selfbot := discord.Selfbot{Token: token}
	gateway := discord.CreateGateway(&selfbot)

	return &Client{&selfbot, gateway}
}

func (client *Client) Connect() error {
	return client.Gateway.Connect()
}

func (client *Client) AddHandler(event string, function any) error {
	return client.Gateway.Handlers.Add(event, function)
}
