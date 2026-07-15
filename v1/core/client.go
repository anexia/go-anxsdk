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

// PermissionGroup returns a permission group api client.
func (c *Client) PermissionGroup() *PermissionGroupClient {
	return newPermissionGroupClient(c.transport)
}

// Locations returns a locations client.
func (c *Client) Locations() *LocationsClient {
	return newLocationsClient(c.transport)
}

// Resources returns a resource api client.
func (c *Client) Resources() *ResourceClient {
	return newResourceClient(c.transport)
}

// Tags returns a tag api client.
func (c *Client) Tags() *TagClient {
	return newTagClient(c.transport)
}
