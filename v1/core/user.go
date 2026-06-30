package core

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

// UserType represents the type of a user.
type UserType string

const (
	// UserTypeUser indicates a user is a normal human user.
	UserTypeUser UserType = "user"

	// UserTypeServiceUser indicates a user is a service account.
	UserTypeServiceUser UserType = "service_user"

	// UserTypeHiddenServiceUser indicates a user is a hidden service account.
	UserTypeHiddenServiceUser UserType = "hidden_service_user"
)

// UserListParams defines the available parameters for the user list endpoint.
type UserListParams struct {
	OrganizationIdentifier string `url:"organization_identifier,omitempty"`
}

type UserSearchParams struct {
	Query                  string `url:"query"`
	OrganizationIdentifier string `url:"organization_identifier,omitempty"`
}

type UserListItem struct {
	Active     bool     `json:"active"`
	Email      string   `json:"email"`
	FirstName  string   `json:"first_name"`
	Identifier string   `json:"identifier"`
	LastName   string   `json:"last_name"`
	Type       UserType `json:"type"`
}

func (i UserListItem) GetID() string {
	return i.Identifier
}

type UserApiAccess struct {
	Basic     bool `json:"basic"`
	Token     bool `json:"token"`
	Signature bool `json:"signature"`
}

type UserGetResponse struct {
	Active                 bool              `json:"active"`
	ApiAccess              UserApiAccess     `json:"api_access"`
	CreatedAt              string            `json:"created_at"`
	Email                  string            `json:"email"`
	FirstName              string            `json:"first_name"`
	PermissionGroups       []common.Resource `json:"permission_groups"`
	Identifier             string            `json:"identifier"`
	Language               string            `json:"language"`
	LastName               string            `json:"last_name"`
	OrganizationName       string            `json:"organization_name"`
	OrganizationIdentifier string            `json:"organization_identifier"`
	Type                   string            `json:"type"`
	Phone                  string            `json:"phone"`
	Mobile                 string            `json:"mobile"`
	Fax                    string            `json:"fax"`
	Position               string            `json:"position"`
	Address                string            `json:"address"`
	Zip                    string            `json:"zip"`
	City                   string            `json:"city"`
	Country                string            `json:"country"`

	// unsure what this is, leaving it commented out
	// ContactResponsibilities []interface{} `json:"contact_responsibilities"`

	TfaEnabled bool `json:"tfa_enabled"`
}

type UserCreateRequest struct {
	LastName               string        `json:"last_name"`
	Type                   UserType      `json:"type"`
	Email                  string        `json:"email"`
	FirstName              string        `json:"firstName,omitempty"`
	Active                 bool          `json:"active"`
	PermissionGroups       []string      `json:"permissionGroups"`
	Language               string        `json:"language,omitempty"`
	Address                string        `json:"address,omitempty"`
	City                   string        `json:"city,omitempty"`
	Country                string        `json:"country,omitempty"`
	Fax                    string        `json:"fax,omitempty"`
	Phone                  string        `json:"phone,omitempty"`
	MobilePhone            string        `json:"mobile_phone,omitempty"`
	Position               string        `json:"position,omitempty"`
	Zip                    string        `json:"zip,omitempty"`
	ApiAccess              UserApiAccess `json:"api_access"`
	OrganizationIdentifier string        `json:"organization_identifier,omitempty"`
	SendPassword           bool          `json:"send_password"`
}

type UserCreateResponse struct {
	Identifier string `json:"identifier"`
}

type UserUpdateRequest struct {
}

type UserUpdateResponse struct {
}

// UserClient is an api client for managing users.
type UserClient struct {
	transport *internal.Transport
}

func newUserClient(transport *internal.Transport) *UserClient {
	return &UserClient{
		transport: transport,
	}
}

// List returns a list of paged users.
func (c *UserClient) List(ctx context.Context, pageParams paging.Params, params UserListParams) (paging.PagedResponse[UserListItem], error) {
	resp := paging.PagedResponse[UserListItem]{}
	err := c.transport.Get(ctx, "/api/core/v1/user.json", &resp, pageParams, params)
	return resp, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for users.
func (c *UserClient) ListPageFetcher(params UserListParams) paging.PageFetcher[UserListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[UserListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// Search returns a paged list of filtered users.
func (c *UserClient) Search(ctx context.Context, pageParams paging.Params, params UserSearchParams) (paging.PagedResponse[UserListItem], error) {
	resp := paging.PagedResponse[UserListItem]{}
	err := c.transport.Get(ctx, "/api/core/v1/user/search.json", &resp, pageParams, params)
	return resp, common.MapTransportError(err)
}

// SearchPageFetcher returns a paging.PageFetcher for users.
func (c *UserClient) SearchPageFetcher(params UserSearchParams) paging.PageFetcher[UserListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[UserListItem], error) {
		return c.Search(ctx, pageParams, params)
	}
}

// Get returns a user by identifier.
func (c *UserClient) Get(ctx context.Context, identifier string) (UserGetResponse, error) {
	resp := UserGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/core/v1/user.json/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}

// Create creates a new user.
func (c *UserClient) Create(ctx context.Context, request UserCreateRequest) (UserCreateResponse, error) {
	resp := UserCreateResponse{}
	err := c.transport.Post(ctx, "/api/core/v1/user.json", request, &resp)
	return resp, common.MapTransportError(err)
}

// Update updates a user.
func (c *UserClient) Update(ctx context.Context, identifier string, request UserUpdateRequest) (UserUpdateResponse, error) {
	resp := UserUpdateResponse{}
	err := c.transport.Put(ctx, fmt.Sprintf("/api/core/v1/user.json/%s", identifier), request, &resp)
	return resp, common.MapTransportError(err)
}

// Delete deletes a user.
func (c *UserClient) Delete(ctx context.Context, identifier string) error {
	err := c.transport.Delete(ctx, fmt.Sprintf("/api/core/v1/user.json/%s", identifier))
	return common.MapTransportError(err)
}
