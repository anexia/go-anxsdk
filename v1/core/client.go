package core

import "github.com/anexia/go-anxsdk/internal"

//go:generate go tool counterfeiter -generate

//counterfeiter:generate -o ../../fake/v1/core . API

// API is the interface implemented by [Client].
type API interface {
	// User returns a user api client.
	User() UserAPI
	// PermissionGroup returns a permission group api client.
	PermissionGroup() PermissionGroupAPI
	// Locations returns a locations client.
	Locations() LocationsAPI
	// Resources returns a resource api client.
	Resources() ResourceAPI
	// Tags returns a tag api client.
	Tags() TagAPI
}

// Client is an anexia core v1 api client.
type Client struct {
	transport *internal.Transport
}

// Check that Client implements the API interface.
var _ API = (*Client)(nil)

// NewClient creates a new v1 core api client.
func NewClient(transport *internal.Transport) *Client {
	return &Client{
		transport: transport,
	}
}

// User returns a user api client.
func (c *Client) User() UserAPI {
	return newUserClient(c.transport)
}

// PermissionGroup returns a permission group api client.
func (c *Client) PermissionGroup() PermissionGroupAPI {
	return newPermissionGroupClient(c.transport)
}

// Locations returns a locations client.
func (c *Client) Locations() LocationsAPI {
	return newLocationsClient(c.transport)
}

// Resources returns a resource api client.
func (c *Client) Resources() ResourceAPI {
	return newResourceClient(c.transport)
}

// Tags returns a tag api client.
func (c *Client) Tags() TagAPI {
	return newTagClient(c.transport)
}
