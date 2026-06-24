package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v2/common"
)

type SyncSource string

const (
	SyncSourceEngine  SyncSource = "engine"
	SyncSourceCluster SyncSource = "cluster"
)

type OS string

const (
	OSFlatcar OS = "Flatcar Linux"
)

const (
	MebiByte = 1024 * 1024
	GibiByte = MebiByte * 1024
)

type CPUPerformanceType string

const (
	CPUPerformanceTypeBestEffort      CPUPerformanceType = "best-effort"
	CPUPerformanceTypeStandard        CPUPerformanceType = "standard"
	CPUPerformanceTypeEnterprise      CPUPerformanceType = "enterprise"
	CPUPerformanceTypePerformance     CPUPerformanceType = "performance"
	CPUPerformanceTypePerformancePlus CPUPerformanceType = "performance-plus"

	CPUPerformanceTypeDefault = CPUPerformanceTypePerformance
)

// State is a type to represent valid nodepool states IDs.
type State string

// All available supported State values.
const (
	NodepoolStateDeployed State = "0"
	NodepoolStateError    State = "1"
)

// NodepoolListParams defines the available parameters for the nodepool list endpoint.
type NodepoolListParams struct {
	Search string `url:"search,omitempty"`

	State               State               `url:"state,omitempty"`
	OperatingSystem     OS                  `url:"operating_system,omitempty"`
	ClusterID           string              `url:"cluster,omitempty"`
	SyncSource          SyncSource          `url:"syncsource,omitempty"`
	CPUPerformanceType  CPUPerformanceType  `url:"cpu_performance_type,omitempty"`
	DiskPerformanceType DiskPerformanceType `url:"disk_performance_type,omitempty"`
	AutoscalerEnabled   bool                `url:"autoscaler_enabled,omitempty"`

	DNSOverrideIPv4 bool   `url:"dns_override_ipv4,omitempty"`
	DNSv4Entry1     string `url:"dns_v4_1,omitempty"`
	DNSv4Entry2     string `url:"dns_v4_2,omitempty"`

	DNSOverrideIPv6 bool   `url:"dns_override_ipv6,omitempty"`
	DNSv6Entry1     string `url:"dns_v6_1,omitempty"`
	DNSv6Entry2     string `url:"dns_v6_2,omitempty"`
}

// NodepoolListItem is an item in the nodepool list response.
type NodepoolListItem common.Resource

// NodepoolGetResponse resource represents the main resource to map to the MachineDeployment in the customer cluster.
type NodepoolGetResponse struct {
	State              common.State[State] `json:"state,omitempty"`
	CustomerIdentifier string              `json:"customer_identifier"`
	ResellerIdentifier string              `json:"reseller_identifier"`
	Identifier         string              `json:"identifier"`
	Name               string              `json:"name"`

	Cluster            common.Resource                         `json:"cluster"`
	SyncSource         common.IDTitleTuple[SyncSource]         `json:"syncsource"`
	Replicas           uint                                    `json:"replicas"`
	CPUs               uint                                    `json:"cpus"`
	CPUType            common.IDTitleTuple[CPUPerformanceType] `json:"cpu_performance_type"`
	MemoryBytes        uint64                                  `json:"memory"`
	OperatingSystem    common.IDTitleTuple[OS]                 `json:"operating_system"`
	AutoscalerEnabled  bool                                    `json:"autoscaler_enabled"`
	AutoscalerMinNodes uint                                    `json:"autoscaler_min_nodes"`
	AutoscalerMaxNodes uint                                    `json:"autoscaler_max_nodes"`

	DiskSizeBytes       uint64                                   `json:"disk_size"`
	DiskPerformanceType common.IDTitleTuple[DiskPerformanceType] `json:"disk_performance_type"`
	AdditionalDisks     []DisksGetResponse                       `json:"additional_disks"`
	Networks            []NetworkGetResponse                     `json:"networks"`

	DNSOverrideIPv4 bool   `json:"dns_override_ipv4"`
	DNSv4Entry1     string `json:"dns_v4_1"`
	DNSv4Entry2     string `json:"dns_v4_2"`

	DNSOverrideIPv6 bool   `json:"dns_override_ipv6"`
	DNSv6Entry1     string `json:"dns_v6_1"`
	DNSv6Entry2     string `json:"dns_v6_2"`

	Taints      string `json:"taints"`
	Labels      string `json:"labels"`
	Annotations string `json:"annotations"`
	SSHPubKeys  string `json:"sshpubkeys"`

	AutomationRules []common.Resource `json:"automation_rules"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type NodepoolUpdateRequest struct {
	// TODO
}

type NodepoolUpdateResponse struct {
	// TODO
}

// NodepoolsClient is an api client for managing Nodepools.
type NodepoolsClient struct {
	environment KubernetesEnv
	transport   *internal.Transport
}

func NewNodepoolsClient(transport *internal.Transport, environment KubernetesEnv) *NodepoolsClient {
	return &NodepoolsClient{
		transport:   transport,
		environment: environment,
	}
}

func (c *NodepoolsClient) endpointRoot() string {
	switch c.environment {
	case KubernetesEnvDevelopment:
		return "/api/kubernetes-dev"
	case KubernetesEnvStaging:
		return "/api/kubernetes-stage"
	case KubernetesEnvProduction:
		return "/api/kubernetes"
	}

	return "/api/kubernetes"
}

// Get returns a single nodepool by its id.
func (c *NodepoolsClient) Get(ctx context.Context, identifier string) (NodepoolGetResponse, error) {
	resp := NodepoolGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("%s/v2/node_pool/%s", c.endpointRoot(), identifier), &resp)
	return resp, common.MapTransportError(err)
}

// Update updates a nodepool.
func (c *NodepoolsClient) Update(ctx context.Context, identifier string, request NodepoolUpdateRequest) (NodepoolUpdateResponse, error) {
	resp := NodepoolUpdateResponse{}
	err := c.transport.Put(ctx, fmt.Sprintf("%s/v2/node_pool/%s", c.endpointRoot(), identifier), request, &resp)
	return resp, common.MapTransportError(err)
}

// List returns a list of paged Nodepools.
func (c *NodepoolsClient) List(ctx context.Context, pageParams paging.Params, params NodepoolListParams) (paging.PagedResponse[NodepoolListItem], error) {
	resp := paging.PagedResponse[NodepoolListItem]{}
	err := c.transport.Get(ctx, fmt.Sprintf("%s/v2/node_pool", c.endpointRoot()), &resp, pageParams, params)
	return resp, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for Nodepools.
func (c *NodepoolsClient) ListPageFetcher(params NodepoolListParams) paging.PageFetcher[NodepoolListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[NodepoolListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFull returns a list of paged Nodepools with all attributes included.
func (c *NodepoolsClient) ListFull(ctx context.Context, pageParams paging.Params, params NodepoolListParams) (paging.PagedResponse[NodepoolGetResponse], error) {
	resp := paging.PagedResponse[NodepoolGetResponse]{}
	wrap := internal.NewAllAttributesWrapper(params)
	err := c.transport.Get(ctx, fmt.Sprintf("%s/v2/node_pool", c.endpointRoot()), &resp, pageParams, wrap)
	return resp, common.MapTransportError(err)
}

// ListFullPageFetcher returns a paging.PageFetcher for Nodepools with all attributes included.
func (c *NodepoolsClient) ListFullPageFetcher(params NodepoolListParams) paging.PageFetcher[NodepoolGetResponse] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[NodepoolGetResponse], error) {
		return c.ListFull(ctx, pageParams, params)
	}
}
