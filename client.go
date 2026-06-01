package go_anx_sdk

import (
	"github.com/kkostial/go-anx-sdk/config"
	"github.com/kkostial/go-anx-sdk/internal"
	v1 "github.com/kkostial/go-anx-sdk/v1"
	v2 "github.com/kkostial/go-anx-sdk/v2"
)

// Client is the entry point to the anexia api.
type Client struct {
	transport *internal.Transport
}

// NewClient creates a new anexia go sdk client with the provided options.
func NewClient(opts ...config.Option) *Client {
	cfg := config.NewConfig(opts...)

	return &Client{
		transport: cfg.CreateTransport(),
	}
}

// V1 returns an entry point to anexia api v1 api clients.
func (c *Client) V1() *v1.Client {
	return v1.NewClient(c.transport)
}

// V2 returns an entry point to anexia api v2 api clients.
func (c *Client) V2() *v2.Client {
	return v2.NewClient(c.transport)
}
