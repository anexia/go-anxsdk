package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v2/common"
)

//revive:disable // self-explanatory constants

type DiskPerformanceType string

const (
	DiskPerformanceTypeSTD1 DiskPerformanceType = "STD1"
	DiskPerformanceTypeSTD2 DiskPerformanceType = "STD2"
	DiskPerformanceTypeSTD3 DiskPerformanceType = "STD3"
	DiskPerformanceTypeSTD4 DiskPerformanceType = "STD4"
	DiskPerformanceTypeSTD5 DiskPerformanceType = "STD5"
	DiskPerformanceTypeENT1 DiskPerformanceType = "ENT1"
	DiskPerformanceTypeENT2 DiskPerformanceType = "ENT2"
	DiskPerformanceTypeENT3 DiskPerformanceType = "ENT3"
	DiskPerformanceTypeENT4 DiskPerformanceType = "ENT4"
	DiskPerformanceTypeENT5 DiskPerformanceType = "ENT5"
	DiskPerformanceTypeENT6 DiskPerformanceType = "ENT6"

	DiskPerformanceTypeDefault = DiskPerformanceTypeENT6
)

//revive:enable

// NodepoolDiskGetResponse represents the disks of a Nodepool.
type NodepoolDiskGetResponse struct {
	Identifier string          `json:"identifier,omitempty"`
	Name       string          `json:"name,omitempty"`
	Nodepool   common.Resource `json:"nodepool,omitempty"`

	SizeBytes       uint64                                   `json:"size_bytes,omitempty"`
	PerformanceType common.IDTitleTuple[DiskPerformanceType] `json:"performance_type,omitempty"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// NodepoolDiskUpdateRequest represents the disks of a Nodepool.
type NodepoolDiskUpdateRequest struct {
	Name       string `json:"name,omitempty"`
	NodepoolID string `json:"nodepool,omitempty"`

	SizeBytes       uint64              `json:"size_bytes,omitempty"`
	PerformanceType DiskPerformanceType `json:"performance_type,omitempty"`
}

// NodepoolDiskListParams defines the available parameters for the disk list endpoint.
type NodepoolDiskListParams struct {
	Search string `url:"search,omitempty"`

	NodepoolID      string              `url:"nodepool,omitempty"`
	PerformanceType DiskPerformanceType `url:"performance_type,omitempty"`
}

// NodepoolDiskListItem is an item in the disk list response.
type NodepoolDiskListItem common.Resource

// NodepoolDisksClient is an api client for managing NodepoolDisks.
type NodepoolDisksClient struct {
	environment Env
	transport   *internal.Transport
}

func newNodepoolDisksClient(transport *internal.Transport, environment Env) *NodepoolDisksClient {
	return &NodepoolDisksClient{
		transport:   transport,
		environment: environment,
	}
}

func (c *NodepoolDisksClient) endpointRoot() string {
	switch c.environment {
	case EnvDevelopment:
		return "/api/kubernetes-dev"
	case EnvStaging:
		return "/api/kubernetes-stage"
	case EnvProduction:
		return "/api/kubernetes"
	}

	return "/api/kubernetes"
}

// Get returns a single disk by its id.
func (c *NodepoolDisksClient) Get(ctx context.Context, identifier string) (NodepoolDiskGetResponse, error) {
	resp := NodepoolDiskGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("%s/v2/node_pool_disk/%s", c.endpointRoot(), identifier), &resp)
	return resp, common.MapTransportError(err)
}

// Update updates a disk.
func (c *NodepoolDisksClient) Update(ctx context.Context, identifier string, request NodepoolDiskUpdateRequest) (NodepoolDiskGetResponse, error) {
	resp := NodepoolDiskGetResponse{}
	err := c.transport.Put(ctx, fmt.Sprintf("%s/v2/node_pool_disk/%s", c.endpointRoot(), identifier), request, &resp)
	return resp, common.MapTransportError(err)
}

// List returns a list of paged Disks.
func (c *NodepoolDisksClient) List(ctx context.Context, pageParams paging.Params, params NodepoolDiskListParams) (paging.PagedResponse[NodepoolDiskListItem], error) {
	resp := paging.PagedResponse[NodepoolDiskListItem]{}
	err := c.transport.Get(ctx, fmt.Sprintf("%s/v2/node_pool_disk", c.endpointRoot()), &resp, pageParams, params)
	return resp, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for Disks.
func (c *NodepoolDisksClient) ListPageFetcher(params NodepoolDiskListParams) paging.PageFetcher[NodepoolDiskListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[NodepoolDiskListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFull returns a list of paged Disks with all attributes included.
func (c *NodepoolDisksClient) ListFull(ctx context.Context, pageParams paging.Params, params NodepoolDiskListParams) (paging.PagedResponse[NodepoolDiskGetResponse], error) {
	resp := paging.PagedResponse[NodepoolDiskGetResponse]{}
	wrap := internal.NewAllAttributesWrapper(params)
	err := c.transport.Get(ctx, fmt.Sprintf("%s/v2/node_pool_disk", c.endpointRoot()), &resp, pageParams, wrap)
	return resp, common.MapTransportError(err)
}

// ListFullPageFetcher returns a paging.PageFetcher for Disks with all attributes included.
func (c *NodepoolDisksClient) ListFullPageFetcher(params NodepoolDiskListParams) paging.PageFetcher[NodepoolDiskGetResponse] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[NodepoolDiskGetResponse], error) {
		return c.ListFull(ctx, pageParams, params)
	}
}
