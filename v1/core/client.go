package core

import "github.com/anexia/go-anxsdk/internal"

// Client is an anexia core v1 api client.
type Client struct {
	transport *internal.Transport
}

// NewClient creates a new v1 core api client.
func NewClient(transport *internal.Transport) *Client {
	return &Client{
		transport: transport,
	}
}

// User returns a user api client.
func (c *Client) User() *UserClient {
	return newUserClient(c.transport)
}
