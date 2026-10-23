package anxsdk

import (
	"github.com/anexia/go-anxsdk/internal"
	v1 "github.com/anexia/go-anxsdk/v1"
	v2 "github.com/anexia/go-anxsdk/v2"
)

// ClientAPI is the interface implemented by [Client].
type ClientAPI interface {
	// V1 returns an entry point to anexia api v1 api clients.
	V1() v1.API
	// V2 returns an entry point to anexia api v2 api clients.
	V2() v2.API
}

// Client is the entry point to the anexia api.
type Client struct {
	transport *internal.Transport
}

// Check that Client implements the ClientAPI interface.
var _ ClientAPI = (*Client)(nil)

// NewClient creates a new anexia go sdk client with the provided options.
func NewClient(opts ...Option) *Client {
	cfg := newConfig(opts...)

	return &Client{
		transport: cfg.createTransport(),
	}
}

// V1 returns an entry point to anexia api v1 api clients.
func (c *Client) V1() v1.API {
	return v1.NewClient(c.transport)
}

// V2 returns an entry point to anexia api v2 api clients.
func (c *Client) V2() v2.API {
	return v2.NewClient(c.transport)
}
