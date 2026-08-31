package core

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

// TagListParams defines the available parameters for the tag list endpoint.
type TagListParams struct {
}

// TagListItem is an item in the tag list response.
type TagListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// TagGetResponse represents a single tag.
type TagGetResponse struct {
	Identifier              string                   `json:"identifier"`
	Name                    string                   `json:"name"`
	OrganizationAssignments []OrganizationAssignment `json:"organisation_assignments"`
}

// OrganizationAssignment represents the assignment of a tag.
type OrganizationAssignment struct {
	Customer Customer        `json:"customer"`
	Service  common.Resource `json:"service"`
}

// TagCreateRequest represents a request to create a new tag.
type TagCreateRequest struct {
	Name              string `json:"name"`
	ServiceIdentifier string `json:"service_identifier"`
}

// TagClient is an api client for managing tags.
type TagClient struct {
	transport *internal.Transport
}

// newTagClient creates a new tag client.
func newTagClient(transport *internal.Transport) *TagClient {
	return &TagClient{
		transport: transport,
	}
}

// List returns a list of paged tags.
func (c *TagClient) List(ctx context.Context, pagingParams paging.Params, params TagListParams) (paging.PagedResponse[TagListItem], error) {
	resp := paging.PagedResponse[TagListItem]{}
	err := c.transport.Get(ctx, "api/core/v1/tags.json", &resp, pagingParams, params)
	return resp, common.MapTransportError(err)
}

// Get returns a single tag by its identifier.
func (c *TagClient) Get(ctx context.Context, identifier string) (TagGetResponse, error) {
	resp := TagGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("api/core/v1/tags.json/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}

// Create creates a new tag.
func (c *TagClient) Create(ctx context.Context, request TagCreateRequest) error {
	err := c.transport.Post(ctx, "api/core/v1/tags.json", request, nil)
	return common.MapTransportError(err)
}

// Delete deletes a tag.
func (c *TagClient) Delete(ctx context.Context, identifier string, serviceIdentifier string) error {
	err := c.transport.Delete(ctx, fmt.Sprintf("api/core/v1/tags.json/%s?service_identifier=%s", identifier, serviceIdentifier))
	return common.MapTransportError(err)
}
