package dgate

import "github.com/LuminalDev/dgate/discord"

type Client struct {
	Selfbot *discord.Selfbot
	Gateway *discord.Gateway
}

func NewClient(token string) *Client {
	s := discord.Selfbot{Token: token}
	g := discord.CreateGateway(&s)
	return &Client{&s, g}
}

func (c *Client) Connect() error {
	return c.Gateway.Connect()
}

func (c *Client) AddHandler(event string, function any) error {
	return c.Gateway.Handlers.Add(event, function)
}
