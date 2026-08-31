package core

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

// PermissionListParams defines the available parameters for the permission list endpoint.
type PermissionListParams struct {
	OrganizationIdentifier string `url:"organization_identifier,omitempty"`
	ServiceIdentifier      string `url:"service_identifier,omitempty"`
}

// PermissionGroupListParams defines the available parameters for the permission group list endpoint.
type PermissionGroupListParams struct {
	OrganizationIdentifier string `url:"organization_identifier,omitempty"`
}

// PermissionGroup is an item in the permission group list response.
type PermissionGroup struct {
	Comment                string            `json:"comment"`
	Identifier             string            `json:"identifier"`
	Name                   string            `json:"name"`
	Permissions            []common.Resource `json:"permissions"`
	Visible                bool              `json:"visible"`
	OrganizationIdentifier string            `url:"organization_identifier"`
}

// PermissionGroupUpdateRequest represents all updatable fields.
type PermissionGroupUpdateRequest struct {
	Comment     *string   `json:"comment,omitempty"`
	Name        *string   `json:"name,omitempty"`
	Permissions *[]string `json:"permissions,omitempty"`
	Visible     *bool     `json:"visible,omitempty"`
}

// PermissionGroupUpdateResponse represents the response of updating a permission group.
type PermissionGroupUpdateResponse struct{}

// PermissionGroupCreateRequest represents the permission group create request.
type PermissionGroupCreateRequest struct {
	Comment                string            `json:"comment"`
	Name                   string            `json:"name"`
	Permissions            []common.Resource `json:"permissions"`
	Visible                bool              `json:"visible"`
	OrganizationIdentifier string            `url:"organization_identifier,omitempty"`
}

// PermissionGroupCreateResponse represents the response of creating a permission group.
type PermissionGroupCreateResponse struct {
	Identifier string `json:"identifier"`
}

// PermissionGroupClient is an api client for managing permission groups.
type PermissionGroupClient struct {
	transport *internal.Transport
}

func newPermissionGroupClient(transport *internal.Transport) *PermissionGroupClient {
	return &PermissionGroupClient{
		transport: transport,
	}
}

// List returns a paged list of permission groups.
func (c *PermissionGroupClient) List(ctx context.Context, pageParams paging.Params, params PermissionGroupListParams) (paging.PagedResponse[PermissionGroup], error) {
	resp := paging.PagedResponse[PermissionGroup]{}
	err := c.transport.Get(ctx, "/api/core/v1/permissiongroup.json", &resp, pageParams, params)
	return resp, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for permission groups.
func (c *PermissionGroupClient) ListPageFetcher(params PermissionGroupListParams) paging.PageFetcher[PermissionGroup] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[PermissionGroup], error) {
		return c.List(ctx, pageParams, params)
	}
}

// Get returns a permission group by identifier.
func (c *PermissionGroupClient) Get(ctx context.Context, identifier string) (PermissionGroup, error) {
	resp := PermissionGroup{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/core/v1/permissiongroup.json/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}

// Create creates a new permission group.
func (c *PermissionGroupClient) Create(ctx context.Context, group PermissionGroupCreateRequest) (PermissionGroupCreateResponse, error) {
	resp := PermissionGroupCreateResponse{}
	err := c.transport.Post(ctx, "/api/core/v1/permissiongroup.json", group, &resp)
	return resp, common.MapTransportError(err)
}

// Update updates a permission group.
func (c *PermissionGroupClient) Update(ctx context.Context, identifier string, request PermissionGroupUpdateRequest) (PermissionGroupUpdateResponse, error) {
	resp := PermissionGroupUpdateResponse{}
	err := c.transport.Put(ctx, fmt.Sprintf("/api/core/v1/permissiongroup.json/%s", identifier), request, &resp)
	return resp, common.MapTransportError(err)
}

// Delete deletes a permission group.
func (c *PermissionGroupClient) Delete(ctx context.Context, identifier string) error {
	err := c.transport.Delete(ctx, fmt.Sprintf("/api/core/v1/permissiongroup.json/%s", identifier))
	return common.MapTransportError(err)
}

// ListPermissions lists all permissions.
func (c *PermissionGroupClient) ListPermissions(ctx context.Context, params PermissionListParams) ([]common.Resource, error) {
	var resp []common.Resource
	err := c.transport.Get(ctx, "/api/core/v1/permissoingroup/permissions.json", &resp, paging.DefaultParams(), params)
	return resp, common.MapTransportError(err)
}
