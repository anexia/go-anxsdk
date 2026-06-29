package vsphere

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

// CPUArchitecture represents the cpu architecture of a vm.
type CPUArchitecture struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// CPUPerformanceType represents the cpu performance characteristics.
type CPUPerformanceType struct {
	ID             string  `json:"id"`
	Architecture   *string `json:"architecture,omitempty"`
	Prioritization string  `json:"prioritization"`
	Limit          float32 `json:"limit"`
	Unit           string  `json:"unit"`
}

// DiskType represents the disk characteristics.
type DiskType struct {
	ID          string `json:"id"`
	StorageType string `json:"storage_type"`
	Bandwidth   int    `json:"bandwidth"`
	IOPs        int    `json:"iops"`
	Latency     int    `json:"latency"`
}

// Location represents a datacenter.
type Location struct {
	Code        string  `json:"code"`
	Country     *string `json:"country"`
	ID          string  `json:"id"`
	Lat         *string `json:"lat"`
	Lon         *string `json:"lon"`
	Name        string  `json:"name"`
	CountryName string  `json:"country_name"`
}

// AvailabilityZone represents a per location sub zone.
type AvailabilityZone struct {
	Identifier     string   `json:"identifier"`
	Name           string   `json:"name"`
	ClusterName    string   `json:"cluster_name"`
	CPUCategories  []string `json:"cpu_categories"`
	DiskCategories []string `json:"disk_categories"`
}

// NicType is the type of a nic.
type NicType string

// ProvisioningProgress represents the current progress of a vm provisioning.
type ProvisioningProgress struct {
	Identifier   string   `json:"identifier"`
	Queued       string   `json:"queued"`
	Progress     int      `json:"progress"`
	VMIdentifier string   `json:"vm_identifier"`
	Errors       []string `json:"errors"`
	Status       string   `json:"status"`
}

// TemplateType is the template type for provisioning vms.
type TemplateType string

const (
	// FromScratchTemplateType indicates that a vm uses the from_scratch template type.
	FromScratchTemplateType TemplateType = "from_scratch"
	// TemplatesTemplateType indicates that a vm uses the templates template type.
	TemplatesTemplateType TemplateType = "templates"
)

// ProvisioningRequest represents a VM provisioning request.
type ProvisioningRequest struct {
	Hostname           string                              `json:"hostname"`
	MemoryMB           *int                                `json:"memory_mb,omitempty"`
	CPUs               *int                                `json:"cpus,omitempty"`
	DiskGB             *int                                `json:"disk_gb,omitempty"`
	DiskType           *string                             `json:"disk_type,omitempty"`
	AdditionalDisks    []ProvisioningRequestAdditionalDisk `json:"additional_disks"`
	CPUPerformanceType *string                             `json:"cpu_performance_type,omitempty"`
	AvailabilityZone   *string                             `json:"availability_zone,omitempty"`
	Sockets            *int                                `json:"sockets,omitempty"`
	Network            []ProvisioningRequestNetwork        `json:"network"`
	VideoMemoryAuto    *bool                               `json:"video_memory_auto,omitempty"`
	VideoMemoryMb      *int                                `json:"video_memory_mb,omitempty"`
	DNS1               *string                             `json:"dns1,omitempty"`
	DNS2               *string                             `json:"dns2,omitempty"`
	DNS3               *string                             `json:"dns3,omitempty"`
	DNS4               *string                             `json:"dns4,omitempty"`
	Password           *string                             `json:"password,omitempty"`
	SSH                *string                             `json:"ssh,omitempty"`
	Script             *string                             `json:"script,omitempty"`
	BootDelay          *int                                `json:"boot_delay,omitempty"`
	EnterBiosSetup     *bool                               `json:"enter_bios_setup,omitempty"`
	CustomName         *string                             `json:"customName,omitempty"`
	VTPMEnabled        *bool                               `json:"vtpm_enabled,omitempty"`
	Firmware           *string                             `json:"firmware,omitempty"`
	OSHostname         *string                             `json:"os_hostname,omitempty"`
}

// ProvisioningResponse represents the response of a provisioning request.
type ProvisioningResponse struct {
	Progress   int      `json:"progress"`
	Errors     []string `json:"errors"`
	Identifier string   `json:"identifier"`
	Queued     bool     `json:"queued"`
}

// ProvisioningRequestAdditionalDisk represents the info needed for an additional disk for vm provisioning.
type ProvisioningRequestAdditionalDisk struct {
	GB   int    `json:"gb"`
	Type string `json:"type"`
}

// ProvisioningRequestNetwork represents the network info for vm provisioning.
type ProvisioningRequestNetwork struct {
	NicType        string   `json:"nic_type"`
	BandwidthLimit int      `json:"bandwidth_limit"`
	Vlan           string   `json:"vlan"`
	Ips            []string `json:"ips"`
}

// ProvisioningClient is an api client for managing vm provisioning.
type ProvisioningClient struct {
	transport *internal.Transport
}

func newProvisioningClient(transport *internal.Transport) *ProvisioningClient {
	return &ProvisioningClient{
		transport: transport,
	}
}

// GetCPUArchitectures returns all available cpu architectures.
func (c *ProvisioningClient) GetCPUArchitectures(ctx context.Context) ([]CPUArchitecture, error) {
	var resp []CPUArchitecture
	err := c.transport.GetSingle(ctx, "/api/vshphere/v1/provisioning/cpu_architecture.json", &resp)
	return resp, common.MapTransportError(err)
}

// GetPerformanceTypes returns all available cpu performance types.
func (c *ProvisioningClient) GetPerformanceTypes(ctx context.Context) ([]CPUPerformanceType, error) {
	var resp []CPUPerformanceType
	err := c.transport.GetSingle(ctx, "/api/vshphere/v1/provisioning/cpu_performance_type.json", &resp)
	return resp, common.MapTransportError(err)
}

// GetDiskTypes returns all available disk types in a location.
func (c *ProvisioningClient) GetDiskTypes(ctx context.Context, locationIdentifier string) ([]DiskType, error) {
	var resp []DiskType
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/vshphere/v1/provisioning/cpu_performance_type.json/%s", locationIdentifier), &resp)
	return resp, common.MapTransportError(err)
}

// ListLocations lists a paged response of locations.
func (c *ProvisioningClient) ListLocations(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[Location], error) {
	resp := paging.PagedResponse[Location]{}
	err := c.transport.Get(ctx, "/api/vshphere/v1/provisioning/locations.json", &resp, pageParams, nil)
	return resp, common.MapTransportError(err)
}

// ListLocationPageFetcher returns a paging.PageFetcher for locations.
func (c *ProvisioningClient) ListLocationPageFetcher() paging.PageFetcher[Location] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[Location], error) {
		return c.ListLocations(ctx, pageParams)
	}
}

// ListAvailabilityZones lists a paged response of availability zones in a location.
func (c *ProvisioningClient) ListAvailabilityZones(ctx context.Context, locationIdentifier string, pageParams paging.Params) (paging.PagedResponse[AvailabilityZone], error) {
	resp := paging.PagedResponse[AvailabilityZone]{}
	err := c.transport.Get(ctx, fmt.Sprintf("/api/vshphere/v1/provisioning/locations.json/%s/availability_zone", locationIdentifier), &resp, pageParams, nil)
	return resp, common.MapTransportError(err)
}

// ListAvailabilityZonesPageFetcher returns a paging.PageFetcher for availability zones in a location.
func (c *ProvisioningClient) ListAvailabilityZonesPageFetcher(locationIdentifier string) paging.PageFetcher[AvailabilityZone] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[AvailabilityZone], error) {
		return c.ListAvailabilityZones(ctx, locationIdentifier, pageParams)
	}
}

// GetNicTypes returns all available nic types.
func (c *ProvisioningClient) GetNicTypes(ctx context.Context) ([]NicType, error) {
	var resp []NicType
	err := c.transport.GetSingle(ctx, "/api/vshphere/v1/provisioning/nic_type.json", &resp)
	return resp, common.MapTransportError(err)
}

// GetProvisioningProgress returns the progress for the specified vm provisioning.
func (c *ProvisioningClient) GetProvisioningProgress(ctx context.Context, progressIdentifier string) (ProvisioningProgress, error) {
	resp := ProvisioningProgress{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/vshphere/v1/provisioning/progress.json/%s", progressIdentifier), &resp)
	return resp, common.MapTransportError(err)
}

// Provision provisions a new vm.
func (c *ProvisioningClient) Provision(
	ctx context.Context,
	locationIdentifier string,
	templateType TemplateType,
	templateIdentifier string,
	request ProvisioningRequest,
) (ProvisioningResponse, error) {
	resp := ProvisioningResponse{}
	err := c.transport.Post(ctx, fmt.Sprintf("/api/vshphere/v1/provisioning/vm.json/%s/%s/%s", locationIdentifier, templateType, templateIdentifier), request, &resp)
	return resp, common.MapTransportError(err)
}
