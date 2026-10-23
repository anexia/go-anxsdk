package anxsdkfake

import (
	"github.com/anexia/go-anxsdk"
	v1fake "github.com/anexia/go-anxsdk/fake/v1"
	corefake "github.com/anexia/go-anxsdk/fake/v1/core"
	ipamfake "github.com/anexia/go-anxsdk/fake/v1/ipam"
	v1lbaasfake "github.com/anexia/go-anxsdk/fake/v1/lbaas"
	vspherefake "github.com/anexia/go-anxsdk/fake/v1/vsphere"
	v2fake "github.com/anexia/go-anxsdk/fake/v2"
	kubernetesfake "github.com/anexia/go-anxsdk/fake/v2/kubernetes"
	v2lbaasfake "github.com/anexia/go-anxsdk/fake/v2/lbaas"
	v1 "github.com/anexia/go-anxsdk/v1"
	v2 "github.com/anexia/go-anxsdk/v2"
)

// Client is a fully pre-wired fake of the [anxsdk.Client]. It implements
// [anxsdk.ClientAPI] and exposes every generated fake via getters, so
// individual clients can be configured before use.
type Client struct {
	v1Fake   *v1fake.FakeAPI
	vlans    *v1fake.FakeVlansAPI
	clusters *v1fake.FakeClustersAPI

	core             *corefake.FakeAPI
	locations        *corefake.FakeLocationsAPI
	users            *corefake.FakeUserAPI
	permissionGroups *corefake.FakePermissionGroupAPI
	resources        *corefake.FakeResourceAPI
	tags             *corefake.FakeTagAPI

	ipam       *ipamfake.FakeAPI
	addresses  *ipamfake.FakeAddressAPI
	aggregates *ipamfake.FakeAggregateAPI
	prefixes   *ipamfake.FakePrefixAPI
	zones      *ipamfake.FakeZoneAPI

	v1Lbaas               *v1lbaasfake.FakeAPI
	backends              *v1lbaasfake.FakeBackendAPI
	backendSwitchingRules *v1lbaasfake.FakeBackendSwitchingRuleAPI
	binds                 *v1lbaasfake.FakeBindAPI
	frontends             *v1lbaasfake.FakeFrontendAPI
	loadBalancers         *v1lbaasfake.FakeLoadBalancerAPI
	acls                  *v1lbaasfake.FakeACLAPI
	maps                  *v1lbaasfake.FakeMapAPI
	rules                 *v1lbaasfake.FakeRuleAPI
	servers               *v1lbaasfake.FakeServerAPI
	ssls                  *v1lbaasfake.FakeSslAPI

	vsphere      *vspherefake.FakeAPI
	powerControl *vspherefake.FakePowerControlAPI
	info         *vspherefake.FakeInfoAPI
	status       *vspherefake.FakeStatusAPI
	provisioning *vspherefake.FakeProvisioningAPI
	search       *vspherefake.FakeSearchAPI

	v2Fake             *v2fake.FakeAPI
	kubernetes         *kubernetesfake.FakeAPI
	v2Clusters         *kubernetesfake.FakeClustersAPI
	v2Nodepools        *kubernetesfake.FakeNodepoolsAPI
	v2NodepoolDisks    *kubernetesfake.FakeNodepoolDisksAPI
	v2NodepoolNetworks *kubernetesfake.FakeNodepoolNetworksAPI

	v2Lbaas                 *v2lbaasfake.FakeAPI
	v2Backends              *v2lbaasfake.FakeBackendAPI
	v2BackendSwitchingRules *v2lbaasfake.FakeBackendSwitchingRuleAPI
	v2Binds                 *v2lbaasfake.FakeBindAPI
	v2Frontends             *v2lbaasfake.FakeFrontendAPI
	v2LoadBalancers         *v2lbaasfake.FakeLoadBalancerAPI
	v2Acls                  *v2lbaasfake.FakeACLAPI
	v2Maps                  *v2lbaasfake.FakeMapAPI
	v2Rules                 *v2lbaasfake.FakeRuleAPI
	v2Servers               *v2lbaasfake.FakeServerAPI
	v2Ssls                  *v2lbaasfake.FakeSslAPI
}

// Check that Client implements the anxsdk.ClientAPI interface.
var _ anxsdk.ClientAPI = (*Client)(nil)

// NewClient creates a new fake client with a fully pre-wired fake tree.
// Every entry client returns the generated fake of the corresponding real
// client, so the whole tree is safe to navigate. Configure individual fakes
// via the getters, e.g. `client.Locations().ListReturns(...)`.
func NewClient() *Client {
	c := &Client{
		v1Fake:                  &v1fake.FakeAPI{},
		vlans:                   &v1fake.FakeVlansAPI{},
		clusters:                &v1fake.FakeClustersAPI{},
		core:                    &corefake.FakeAPI{},
		locations:               &corefake.FakeLocationsAPI{},
		users:                   &corefake.FakeUserAPI{},
		permissionGroups:        &corefake.FakePermissionGroupAPI{},
		resources:               &corefake.FakeResourceAPI{},
		tags:                    &corefake.FakeTagAPI{},
		ipam:                    &ipamfake.FakeAPI{},
		addresses:               &ipamfake.FakeAddressAPI{},
		aggregates:              &ipamfake.FakeAggregateAPI{},
		prefixes:                &ipamfake.FakePrefixAPI{},
		zones:                   &ipamfake.FakeZoneAPI{},
		v1Lbaas:                 &v1lbaasfake.FakeAPI{},
		backends:                &v1lbaasfake.FakeBackendAPI{},
		backendSwitchingRules:   &v1lbaasfake.FakeBackendSwitchingRuleAPI{},
		binds:                   &v1lbaasfake.FakeBindAPI{},
		frontends:               &v1lbaasfake.FakeFrontendAPI{},
		loadBalancers:           &v1lbaasfake.FakeLoadBalancerAPI{},
		acls:                    &v1lbaasfake.FakeACLAPI{},
		maps:                    &v1lbaasfake.FakeMapAPI{},
		rules:                   &v1lbaasfake.FakeRuleAPI{},
		servers:                 &v1lbaasfake.FakeServerAPI{},
		ssls:                    &v1lbaasfake.FakeSslAPI{},
		vsphere:                 &vspherefake.FakeAPI{},
		powerControl:            &vspherefake.FakePowerControlAPI{},
		info:                    &vspherefake.FakeInfoAPI{},
		status:                  &vspherefake.FakeStatusAPI{},
		provisioning:            &vspherefake.FakeProvisioningAPI{},
		search:                  &vspherefake.FakeSearchAPI{},
		v2Fake:                  &v2fake.FakeAPI{},
		kubernetes:              &kubernetesfake.FakeAPI{},
		v2Clusters:              &kubernetesfake.FakeClustersAPI{},
		v2Nodepools:             &kubernetesfake.FakeNodepoolsAPI{},
		v2NodepoolDisks:         &kubernetesfake.FakeNodepoolDisksAPI{},
		v2NodepoolNetworks:      &kubernetesfake.FakeNodepoolNetworksAPI{},
		v2Lbaas:                 &v2lbaasfake.FakeAPI{},
		v2Backends:              &v2lbaasfake.FakeBackendAPI{},
		v2BackendSwitchingRules: &v2lbaasfake.FakeBackendSwitchingRuleAPI{},
		v2Binds:                 &v2lbaasfake.FakeBindAPI{},
		v2Frontends:             &v2lbaasfake.FakeFrontendAPI{},
		v2LoadBalancers:         &v2lbaasfake.FakeLoadBalancerAPI{},
		v2Acls:                  &v2lbaasfake.FakeACLAPI{},
		v2Maps:                  &v2lbaasfake.FakeMapAPI{},
		v2Rules:                 &v2lbaasfake.FakeRuleAPI{},
		v2Servers:               &v2lbaasfake.FakeServerAPI{},
		v2Ssls:                  &v2lbaasfake.FakeSslAPI{},
	}

	c.wire()

	return c
}

// wire connects every entry fake to its sub-fakes so that navigating the tree
// never returns nil.
func (c *Client) wire() {
	c.v1Fake.CoreReturns(c.core)
	c.v1Fake.IpamReturns(c.ipam)
	c.v1Fake.LBaaSReturns(c.v1Lbaas)
	c.v1Fake.VSphereReturns(c.vsphere)
	c.v1Fake.VlansReturns(c.vlans)
	c.v1Fake.ClustersReturns(c.clusters)
	c.v1Fake.StageClustersReturns(c.clusters)
	c.v1Fake.DevClustersReturns(c.clusters)

	c.core.UserReturns(c.users)
	c.core.PermissionGroupReturns(c.permissionGroups)
	c.core.LocationsReturns(c.locations)
	c.core.ResourcesReturns(c.resources)
	c.core.TagsReturns(c.tags)

	c.ipam.AddressesReturns(c.addresses)
	c.ipam.AggregatesReturns(c.aggregates)
	c.ipam.PrefixesReturns(c.prefixes)
	c.ipam.ZonesReturns(c.zones)

	c.v1Lbaas.BackendsReturns(c.backends)
	c.v1Lbaas.BackendSwitchingRulesReturns(c.backendSwitchingRules)
	c.v1Lbaas.BindsReturns(c.binds)
	c.v1Lbaas.FrontendsReturns(c.frontends)
	c.v1Lbaas.LoadBalancersReturns(c.loadBalancers)
	c.v1Lbaas.AclsReturns(c.acls)
	c.v1Lbaas.MapsReturns(c.maps)
	c.v1Lbaas.RulesReturns(c.rules)
	c.v1Lbaas.ServersReturns(c.servers)
	c.v1Lbaas.SslsReturns(c.ssls)

	c.vsphere.PowerControlReturns(c.powerControl)
	c.vsphere.InfoReturns(c.info)
	c.vsphere.StatusReturns(c.status)
	c.vsphere.ProvisioningReturns(c.provisioning)
	c.vsphere.SearchReturns(c.search)

	c.v2Fake.KubernetesReturns(c.kubernetes)
	c.v2Fake.KubernetesByEnvReturns(c.kubernetes)
	c.v2Fake.KubernetesStageReturns(c.kubernetes)
	c.v2Fake.KubernetesDevReturns(c.kubernetes)
	c.v2Fake.LBaaSReturns(c.v2Lbaas)

	c.kubernetes.ClustersReturns(c.v2Clusters)
	c.kubernetes.NodepoolsReturns(c.v2Nodepools)
	c.kubernetes.NodepoolDisksReturns(c.v2NodepoolDisks)
	c.kubernetes.NodepoolNetworksReturns(c.v2NodepoolNetworks)

	c.v2Lbaas.BackendsReturns(c.v2Backends)
	c.v2Lbaas.BackendSwitchingRulesReturns(c.v2BackendSwitchingRules)
	c.v2Lbaas.BindsReturns(c.v2Binds)
	c.v2Lbaas.FrontendsReturns(c.v2Frontends)
	c.v2Lbaas.LoadBalancersReturns(c.v2LoadBalancers)
	c.v2Lbaas.AclsReturns(c.v2Acls)
	c.v2Lbaas.MapsReturns(c.v2Maps)
	c.v2Lbaas.RulesReturns(c.v2Rules)
	c.v2Lbaas.ServersReturns(c.v2Servers)
	c.v2Lbaas.SslsReturns(c.v2Ssls)
}

// V1 returns the fake v1 api client.
func (c *Client) V1() v1.API {
	return c.v1Fake
}

// V2 returns the fake v2 api client.
func (c *Client) V2() v2.API {
	return c.v2Fake
}

// V1Fake returns the fake v1 entry client.
func (c *Client) V1Fake() *v1fake.FakeAPI {
	return c.v1Fake
}

// V2Fake returns the fake v2 entry client.
func (c *Client) V2Fake() *v2fake.FakeAPI {
	return c.v2Fake
}

// Vlans returns the fake vlans client.
func (c *Client) Vlans() *v1fake.FakeVlansAPI {
	return c.vlans
}

// Clusters returns the fake clusters client.
func (c *Client) Clusters() *v1fake.FakeClustersAPI {
	return c.clusters
}

// StageClusters returns the fake stage clusters client.
func (c *Client) StageClusters() *v1fake.FakeClustersAPI {
	return c.clusters
}

// DevClusters returns the fake dev clusters client.
func (c *Client) DevClusters() *v1fake.FakeClustersAPI {
	return c.clusters
}

// Core returns the fake core client.
func (c *Client) Core() *corefake.FakeAPI {
	return c.core
}

// Locations returns the fake locations client.
func (c *Client) Locations() *corefake.FakeLocationsAPI {
	return c.locations
}

// User returns the fake user client.
func (c *Client) User() *corefake.FakeUserAPI {
	return c.users
}

// PermissionGroup returns the fake permission group client.
func (c *Client) PermissionGroup() *corefake.FakePermissionGroupAPI {
	return c.permissionGroups
}

// Resources returns the fake resource client.
func (c *Client) Resources() *corefake.FakeResourceAPI {
	return c.resources
}

// Tags returns the fake tag client.
func (c *Client) Tags() *corefake.FakeTagAPI {
	return c.tags
}

// Ipam returns the fake ipam client.
func (c *Client) Ipam() *ipamfake.FakeAPI {
	return c.ipam
}

// Addresses returns the fake address client.
func (c *Client) Addresses() *ipamfake.FakeAddressAPI {
	return c.addresses
}

// Aggregates returns the fake aggregate client.
func (c *Client) Aggregates() *ipamfake.FakeAggregateAPI {
	return c.aggregates
}

// Prefixes returns the fake prefix client.
func (c *Client) Prefixes() *ipamfake.FakePrefixAPI {
	return c.prefixes
}

// Zones returns the fake zone client.
func (c *Client) Zones() *ipamfake.FakeZoneAPI {
	return c.zones
}

// LBaaS returns the fake v1 lbaas client.
func (c *Client) LBaaS() *v1lbaasfake.FakeAPI {
	return c.v1Lbaas
}

// Backends returns the fake backend client.
func (c *Client) Backends() *v1lbaasfake.FakeBackendAPI {
	return c.backends
}

// BackendSwitchingRules returns the fake backend switching rule client.
func (c *Client) BackendSwitchingRules() *v1lbaasfake.FakeBackendSwitchingRuleAPI {
	return c.backendSwitchingRules
}

// Binds returns the fake bind client.
func (c *Client) Binds() *v1lbaasfake.FakeBindAPI {
	return c.binds
}

// Frontends returns the fake frontend client.
func (c *Client) Frontends() *v1lbaasfake.FakeFrontendAPI {
	return c.frontends
}

// LoadBalancers returns the fake load balancer client.
func (c *Client) LoadBalancers() *v1lbaasfake.FakeLoadBalancerAPI {
	return c.loadBalancers
}

// Acls returns the fake acl client.
func (c *Client) Acls() *v1lbaasfake.FakeACLAPI {
	return c.acls
}

// Maps returns the fake map client.
func (c *Client) Maps() *v1lbaasfake.FakeMapAPI {
	return c.maps
}

// Rules returns the fake rule client.
func (c *Client) Rules() *v1lbaasfake.FakeRuleAPI {
	return c.rules
}

// Servers returns the fake server client.
func (c *Client) Servers() *v1lbaasfake.FakeServerAPI {
	return c.servers
}

// Ssls returns the fake ssl client.
func (c *Client) Ssls() *v1lbaasfake.FakeSslAPI {
	return c.ssls
}

// VSphere returns the fake vsphere client.
func (c *Client) VSphere() *vspherefake.FakeAPI {
	return c.vsphere
}

// PowerControl returns the fake power control client.
func (c *Client) PowerControl() *vspherefake.FakePowerControlAPI {
	return c.powerControl
}

// Info returns the fake info client.
func (c *Client) Info() *vspherefake.FakeInfoAPI {
	return c.info
}

// Status returns the fake status client.
func (c *Client) Status() *vspherefake.FakeStatusAPI {
	return c.status
}

// Provisioning returns the fake provisioning client.
func (c *Client) Provisioning() *vspherefake.FakeProvisioningAPI {
	return c.provisioning
}

// Search returns the fake search client.
func (c *Client) Search() *vspherefake.FakeSearchAPI {
	return c.search
}

// Kubernetes returns the fake kubernetes client.
func (c *Client) Kubernetes() *kubernetesfake.FakeAPI {
	return c.kubernetes
}

// KubernetesStage returns the fake stage kubernetes client.
func (c *Client) KubernetesStage() *kubernetesfake.FakeAPI {
	return c.kubernetes
}

// KubernetesDev returns the fake dev kubernetes client.
func (c *Client) KubernetesDev() *kubernetesfake.FakeAPI {
	return c.kubernetes
}

// V2Clusters returns the fake v2 clusters client.
func (c *Client) V2Clusters() *kubernetesfake.FakeClustersAPI {
	return c.v2Clusters
}

// V2Nodepools returns the fake v2 nodepools client.
func (c *Client) V2Nodepools() *kubernetesfake.FakeNodepoolsAPI {
	return c.v2Nodepools
}

// V2NodepoolDisks returns the fake v2 nodepool disks client.
func (c *Client) V2NodepoolDisks() *kubernetesfake.FakeNodepoolDisksAPI {
	return c.v2NodepoolDisks
}

// V2NodepoolNetworks returns the fake v2 nodepool networks client.
func (c *Client) V2NodepoolNetworks() *kubernetesfake.FakeNodepoolNetworksAPI {
	return c.v2NodepoolNetworks
}

// V2LBaaS returns the fake v2 lbaas client.
func (c *Client) V2LBaaS() *v2lbaasfake.FakeAPI {
	return c.v2Lbaas
}

// V2Backends returns the fake v2 backend client.
func (c *Client) V2Backends() *v2lbaasfake.FakeBackendAPI {
	return c.v2Backends
}

// V2BackendSwitchingRules returns the fake v2 backend switching rule client.
func (c *Client) V2BackendSwitchingRules() *v2lbaasfake.FakeBackendSwitchingRuleAPI {
	return c.v2BackendSwitchingRules
}

// V2Binds returns the fake v2 bind client.
func (c *Client) V2Binds() *v2lbaasfake.FakeBindAPI {
	return c.v2Binds
}

// V2Frontends returns the fake v2 frontend client.
func (c *Client) V2Frontends() *v2lbaasfake.FakeFrontendAPI {
	return c.v2Frontends
}

// V2LoadBalancers returns the fake v2 load balancer client.
func (c *Client) V2LoadBalancers() *v2lbaasfake.FakeLoadBalancerAPI {
	return c.v2LoadBalancers
}

// V2Acls returns the fake v2 acl client.
func (c *Client) V2Acls() *v2lbaasfake.FakeACLAPI {
	return c.v2Acls
}

// V2Maps returns the fake v2 map client.
func (c *Client) V2Maps() *v2lbaasfake.FakeMapAPI {
	return c.v2Maps
}

// V2Rules returns the fake v2 rule client.
func (c *Client) V2Rules() *v2lbaasfake.FakeRuleAPI {
	return c.v2Rules
}

// V2Servers returns the fake v2 server client.
func (c *Client) V2Servers() *v2lbaasfake.FakeServerAPI {
	return c.v2Servers
}

// V2Ssls returns the fake v2 ssl client.
func (c *Client) V2Ssls() *v2lbaasfake.FakeSslAPI {
	return c.v2Ssls
}
