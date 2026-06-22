package v2

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v2/common"
)

// KubernetesEnv is used to specify the environment to use.
type KubernetesEnv string

const (
	// KubernetesEnvProduction specifies the production environment.
	KubernetesEnvProduction KubernetesEnv = "prod"
	// KubernetesEnvStaging specifies the staging environment.
	KubernetesEnvStaging KubernetesEnv = "stage"
	// KubernetesEnvDevelopment specifies the development environment.
	KubernetesEnvDevelopment KubernetesEnv = "dev"
)

// ClusterState is a type to represent valid cluster states IDs.
type ClusterState string

// All available supported ClusterState values.
const (
	ClusterStateOk                     ClusterState = "0"
	ClusterStateError                  ClusterState = "1"
	CLusterStatePending                ClusterState = "2"
	ClusterStateAnxDevNoGa             ClusterState = "3"
	ClusterStateWaitingForNetworks     ClusterState = "4"
	ClusterStateWaitingForServiceVMs   ClusterState = "5"
	ClusterStateWaitingForControlPlane ClusterState = "6"
	ClusterStateUpdatingControlPlane   ClusterState = "7"
	ClusterStateUpdatingNodes          ClusterState = "8"
)

// ClusterVersion is a type to represent valid cluster versions.
type ClusterVersion string

// All available supported ClusterVersion values.
const (
	ClusterVersion1_32 ClusterVersion = "1.32"
	ClusterVersion1_33 ClusterVersion = "1.33"
	ClusterVersion1_34 ClusterVersion = "1.34"
	ClusterVersion1_35 ClusterVersion = "1.35"
)

// ClusterCniPlugin is a type to represent valid cni plugins.
type ClusterCniPlugin string

// All available supported ClusterCniPlugin values.
const (
	ClusterCniPluginCanal ClusterCniPlugin = "Canal"
)

// ClusterListParams defines the available parameters for the cluster list endpoint.
type ClusterListParams struct {
	Reseller           string `url:"reseller,omitempty"`
	Customer           string `url:"customer,omitempty"`
	Name               string `url:"name,omitempty"`
	State              string `url:"state,omitempty"`
	NeedsServiceVMs    bool   `url:"need_services_vms,omitempty"`
	Location           string `url:"location,omitempty"`
	CNIPlugin          string `url:"cni_plugin,omitempty"`
	Autoscaling        bool   `url:"autoscaling,omitempty"`
	ExternalIPFamilies string `url:"external_ip_families,omitempty"`
}

// ClusterListItem is an item in the cluster list response.
type ClusterListItem common.Resource

// ClusterGetResponse represents the response of the cluster get endpoint.
type ClusterGetResponse struct {
	CustomerIdentifier             string                                `json:"customer_identifier"`
	ResellerIdentifier             string                                `json:"reseller_identifier"`
	Identifier                     string                                `json:"identifier"`
	Name                           string                                `json:"name"`
	State                          common.State[ClusterState]            `json:"state"`
	Location                       common.Resource                       `json:"location"`
	Version                        common.IDTitleTuple[ClusterVersion]   `json:"version"`
	PatchVersion                   string                                `json:"patch_version"`
	Kubeconfig                     string                                `json:"kubeconfig"`
	Autoscaling                    bool                                  `json:"autoscaling"`
	EnablePersistentStorage        bool                                  `json:"enable_persistent_storage"`
	CniPlugin                      common.IDTitleTuple[ClusterCniPlugin] `json:"cni_plugin"`
	APIServerAllowlist             string                                `json:"apiserver_allowlist"`
	BackendName                    common.IDTitleTuple[string]           `json:"backend_name"`
	Backend                        string                                `json:"backend"`
	MaintenanceWindowStartTime     string                                `json:"maintenance_window_start_time"`
	MaintenanceWindowDuration      string                                `json:"maintenance_window_duration"`
	ServiceUser                    common.Resource                       `json:"service_user"`
	ManageInternalIpv4Prefix       bool                                  `json:"manage_internal_ipv4_prefix"`
	InternalIpv4Prefix             common.Resource                       `json:"internal_ipv4_prefix"`
	NeedsServiceVms                bool                                  `json:"needs_service_vms"`
	EnableNatGateways              bool                                  `json:"enable_nat_gateways"`
	EnableLbaas                    bool                                  `json:"enable_lbaas"`
	ExternalIPFamilies             common.IDTitleTuple[string]           `json:"external_ip_families"`
	ManageExternalIpv4Prefix       bool                                  `json:"manage_external_ipv4_prefix"`
	ExternalIpv4Prefix             common.Resource                       `json:"external_ipv4_prefix"`
	ManageExternalIpv6Prefix       bool                                  `json:"manage_external_ipv6_prefix"`
	ExternalIpv6Prefix             common.Resource                       `json:"external_ipv6_prefix"`
	ServiceVM01                    common.Resource                       `json:"service_vm_01"`
	ServiceVM02                    common.Resource                       `json:"service_vm_02"`
	ServiceVM01InternalIpv4Address common.Resource                       `json:"service_vm_01_internal_ipv4_address"`
	ServiceVM02InternalIpv4Address common.Resource                       `json:"service_vm_02_internal_ipv4_address"`
	ServiceVM01ExternalIpv4Address common.Resource                       `json:"service_vm_01_external_ipv4_address"`
	ServiceVM02ExternalIpv4Address common.Resource                       `json:"service_vm_02_external_ipv4_address"`
	ServiceVM01ExternalIpv6Address common.Resource                       `json:"service_vm_01_external_ipv6_address"`
	ServiceVM02ExternalIpv6Address common.Resource                       `json:"service_vm_02_external_ipv6_address"`
	ServiceLb01                    common.Resource                       `json:"service_lb_01"`
	ServiceLb02                    common.Resource                       `json:"service_lb_02"`
	ExternalIpv4Vip                common.Resource                       `json:"external_ipv4_vip"`
	ExternalIpv6Vip                common.Resource                       `json:"external_ipv6_vip"`
	KkpAPILbaasBackend01           common.Resource                       `json:"kkp_api_lbaas_backend_01"`
	KkpAPILbaasBackend02           common.Resource                       `json:"kkp_api_lbaas_backend_02"`
	KKPVpnLbaasBackend01           common.Resource                       `json:"kkp_vpn_lbaas_backend_01"`
	KKPVpnLbaasBackend02           common.Resource                       `json:"kkp_vpn_lbaas_backend_02"`
	StorageServerInterfaceAddress  common.Resource                       `json:"storage_server_interface_address"`
	KKPProjectID                   string                                `json:"kkp_project_id"`
	KKPClusterID                   string                                `json:"kkp_cluster_id"`
	EnableOidcAuthentication       bool                                  `json:"enable_oidc_authentication"`
	OIDCClientID                   string                                `json:"oidc_client_id"`
	OIDCIssuerURL                  string                                `json:"oidc_issuer_url"`
	OidcGroupsClaim                string                                `json:"oidc_groups_claim"`
	OidcUsernameClaim              string                                `json:"oidc_username_claim"`
	OidcExtraScopes                string                                `json:"oidc_extra_scopes"`
	OidcGroupsPrefix               string                                `json:"oidc_groups_prefix"`
	OidcRequiredClaim              string                                `json:"oidc_required_claim"`
	OidcUsernamePrefix             string                                `json:"oidc_username_prefix"`
	AutomationRules                []common.Resource                     `json:"automation_rules"`
}

// ClusterUpdateRequest represents all changes made to a cluster during an update request.
type ClusterUpdateRequest struct {
	Name                       *string           `json:"name,omitempty"`
	State                      *ClusterState     `json:"state,omitempty"`
	LocationIdentifier         *string           `json:"location,omitempty"`
	Version                    *ClusterVersion   `json:"version,omitempty"`
	PatchVersion               *string           `json:"patch_version,omitempty"`
	KubeConfig                 *string           `json:"kubeconfig,omitempty"`
	Autoscaling                *bool             `json:"autoscaling,omitempty"`
	EnablePersistentStorage    *bool             `json:"enable_persistent_storage,omitempty"`
	CniPlugin                  *ClusterCniPlugin `json:"cni_plugin,omitempty"`
	APIServerAllowList         *string           `json:"api_server_allow_list,omitempty"`
	BackendName                *string           `json:"backend_name,omitempty"`
	Backend                    *string           `json:"backend,omitempty"`
	MaintenanceWindowStartTime *string           `json:"maintenance_window_start_time,omitempty"`
	MaintenanceWindowDuration  *string           `json:"maintenance_window_duration,omitempty"`
	ServiceUserIdentifier      *string           `json:"service_user,omitempty"`
	// manage_internal_ipv4_prefix
}

// ClusterUpdateResponse represents the response of the cluster update endpoint.
type ClusterUpdateResponse struct {
	CustomerIdentifier             string                     `json:"customer_identifier"`
	ResellerIdentifier             string                     `json:"reseller_identifier"`
	Identifier                     string                     `json:"identifier"`
	Name                           string                     `json:"name"`
	State                          common.State[ClusterState] `json:"state"`
	Location                       common.Resource            `json:"location"`
	Version                        string                     `json:"version"`
	PatchVersion                   any                        `json:"patch_version"`
	Kubeconfig                     string                     `json:"kubeconfig"`
	Autoscaling                    bool                       `json:"autoscaling"`
	EnablePersistentStorage        bool                       `json:"enable_persistent_storage"`
	CniPlugin                      string                     `json:"cni_plugin"`
	ApiserverAllowlist             any                        `json:"apiserver_allowlist"`
	BackendName                    string                     `json:"backend_name"`
	Backend                        string                     `json:"backend"`
	MaintenanceWindowStartTime     string                     `json:"maintenance_window_start_time"`
	MaintenanceWindowDuration      string                     `json:"maintenance_window_duration"`
	ServiceUser                    common.Resource            `json:"service_user"`
	ManageInternalIpv4Prefix       bool                       `json:"manage_internal_ipv4_prefix"`
	InternalIpv4Prefix             common.Resource            `json:"internal_ipv4_prefix"`
	NeedsServiceVms                bool                       `json:"needs_service_vms"`
	EnableNatGateways              bool                       `json:"enable_nat_gateways"`
	EnableLbaas                    bool                       `json:"enable_lbaas"`
	ExternalIPFamilies             string                     `json:"external_ip_families"`
	ManageExternalIpv4Prefix       bool                       `json:"manage_external_ipv4_prefix"`
	ExternalIpv4Prefix             common.Resource            `json:"external_ipv4_prefix"`
	ManageExternalIpv6Prefix       bool                       `json:"manage_external_ipv6_prefix"`
	ExternalIpv6Prefix             common.Resource            `json:"external_ipv6_prefix"`
	ServiceVM01                    common.Resource            `json:"service_vm_01"`
	ServiceVM02                    common.Resource            `json:"service_vm_02"`
	ServiceVM01InternalIpv4Address common.Resource            `json:"service_vm_01_internal_ipv4_address"`
	ServiceVM02InternalIpv4Address common.Resource            `json:"service_vm_02_internal_ipv4_address"`
	ServiceVM01ExternalIpv4Address common.Resource            `json:"service_vm_01_external_ipv4_address"`
	ServiceVM02ExternalIpv4Address common.Resource            `json:"service_vm_02_external_ipv4_address"`
	ServiceVM01ExternalIpv6Address common.Resource            `json:"service_vm_01_external_ipv6_address"`
	ServiceVM02ExternalIpv6Address common.Resource            `json:"service_vm_02_external_ipv6_address"`
	ServiceLb01                    common.Resource            `json:"service_lb_01"`
	ServiceLb02                    common.Resource            `json:"service_lb_02"`
	ExternalIpv4Vip                common.Resource            `json:"external_ipv4_vip"`
	ExternalIpv6Vip                common.Resource            `json:"external_ipv6_vip"`
	KKPApiLbaasBackend01           common.Resource            `json:"kkp_api_lbaas_backend_01"`
	KKPApiLbaasBackend02           common.Resource            `json:"kkp_api_lbaas_backend_02"`
	KKPVpnLbaasBackend01           common.Resource            `json:"kkp_vpn_lbaas_backend_01"`
	KKPVpnLbaasBackend02           common.Resource            `json:"kkp_vpn_lbaas_backend_02"`
	StorageServerInterfaceAddress  any                        `json:"storage_server_interface_address"`
	KKPProjectID                   string                     `json:"kkp_project_id"`
	KKPClusterID                   string                     `json:"kkp_cluster_id"`
	InternalVlan                   any                        `json:"internal_vlan"`
	ExternalVlan                   any                        `json:"external_vlan"`
	EnableOidcAuthentication       bool                       `json:"enable_oidc_authentication"`
	OIDCClientID                   any                        `json:"oidc_client_id"`
	OIDCIssuerURL                  any                        `json:"oidc_issuer_url"`
	OIDCGroupsClaim                string                     `json:"oidc_groups_claim"`
	OIDCUsernameClaim              string                     `json:"oidc_username_claim"`
	OIDCExtraScopes                any                        `json:"oidc_extra_scopes"`
	OIDCGroupsPrefix               any                        `json:"oidc_groups_prefix"`
	OIDCRequiredClaim              any                        `json:"oidc_required_claim"`
	OIDCUsernamePrefix             any                        `json:"oidc_username_prefix"`
	AutomationRules                []common.Resource          `json:"automation_rules"`
}

// ClustersClient is an api client for managing clusters.
type ClustersClient struct {
	environment KubernetesEnv
	transport   *internal.Transport
}

func newClustersClient(
	transport *internal.Transport,
	environment KubernetesEnv,
) *ClustersClient {
	return &ClustersClient{
		transport:   transport,
		environment: environment,
	}
}

func (c *ClustersClient) endpointRoot() string {
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

// Get returns a single cluster by its id.
func (c *ClustersClient) Get(ctx context.Context, identifier string) (ClusterGetResponse, error) {
	resp := ClusterGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("%s/v2/cluster/%s", c.endpointRoot(), identifier), &resp)
	return resp, common.MapTransportError(err)
}

// Update updates a cluster.
func (c *ClustersClient) Update(ctx context.Context, identifier string, request ClusterUpdateRequest) (ClusterUpdateResponse, error) {
	resp := ClusterUpdateResponse{}
	err := c.transport.Put(ctx, fmt.Sprintf("%s/v2/cluster/%s", c.endpointRoot(), identifier), request, &resp)
	return resp, common.MapTransportError(err)
}

// List returns a list of paged clusters.
func (c *ClustersClient) List(ctx context.Context, pageParams paging.Params, params ClusterListParams) (paging.PagedResponse[ClusterListItem], error) {
	resp := paging.PagedResponse[ClusterListItem]{}
	err := c.transport.Get(ctx, fmt.Sprintf("%s/v2/cluster", c.endpointRoot()), &resp, pageParams, params)
	return resp, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for clusters.
func (c *ClustersClient) ListPageFetcher(params ClusterListParams) paging.PageFetcher[ClusterListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[ClusterListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFull returns a list of paged clusters with all attributes included.
func (c *ClustersClient) ListFull(ctx context.Context, pageParams paging.Params, params ClusterListParams) (paging.PagedResponse[ClusterGetResponse], error) {
	resp := paging.PagedResponse[ClusterGetResponse]{}
	wrap := internal.NewAllAttributesWrapper(params)
	err := c.transport.Get(ctx, fmt.Sprintf("%s/v2/cluster", c.endpointRoot()), &resp, pageParams, wrap)
	return resp, common.MapTransportError(err)
}

// ListFullPageFetcher returns a paging.PageFetcher for clusters with all attributes included.
func (c *ClustersClient) ListFullPageFetcher(params ClusterListParams) paging.PageFetcher[ClusterGetResponse] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[ClusterGetResponse], error) {
		return c.ListFull(ctx, pageParams, params)
	}
}

// RequestKubeconfig triggers the automation endpoint request_kubeconfig.
func (c *ClustersClient) RequestKubeconfig(ctx context.Context, identifier string) error {
	err := c.transport.Post(ctx, fmt.Sprintf("%s/v2/cluster/%s/trigger/request_kubeconfig", c.endpointRoot(), identifier), nil, nil)
	return common.MapTransportError(err)
}

// RemoveKubeconfig triggers the automation endpoint remove_kubeconfig.
func (c *ClustersClient) RemoveKubeconfig(ctx context.Context, identifier string) error {
	err := c.transport.Post(ctx, fmt.Sprintf("%s/v2/cluster/%s/trigger/remove_kubeconfig", c.endpointRoot(), identifier), nil, nil)
	return common.MapTransportError(err)
}
