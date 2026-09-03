package vsphere

import (
	"context"
	"fmt"
	"time"

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

// AvailabilityZone represents a per-location sub zone.
type AvailabilityZone struct {
	Identifier     string   `json:"identifier"`
	Name           string   `json:"name"`
	ClusterName    string   `json:"cluster_name"`
	CPUCategories  []string `json:"cpu_categories"`
	DiskCategories []string `json:"disk_categories"`
}

// NicType is the type of network interface card.
type NicType string

// ProvisioningProgress represents the current progress of a vm provisioning.
type ProvisioningProgress struct {
	TaskIdentifier string             `json:"identifier"`
	Queued         bool               `json:"queued"`
	Progress       int                `json:"progress"`
	VMIdentifier   string             `json:"vm_identifier"`
	Errors         []string           `json:"errors"`
	Status         ProvisioningStatus `json:"status"`
}

// ProvisioningStatus specifies the status of the provisioning request.
type ProvisioningStatus string

const (
	// ProvisioningStatusFailed indicates that the provisioning failed.
	ProvisioningStatusFailed ProvisioningStatus = "-1"
	// ProvisioningStatusSuccess indicates that the provisioning succeeded.
	ProvisioningStatusSuccess ProvisioningStatus = "1"
	// ProvisioningStatusInProgress indicates that the provisioning is still ongoing.
	ProvisioningStatusInProgress ProvisioningStatus = "2"
	// ProvisioningStatusCancelled indicates that the provisioning has been cancelled.
	ProvisioningStatusCancelled ProvisioningStatus = "3"
)

// TemplateType is the template type for provisioning vms.
type TemplateType string

const (
	// TemplateTypeFromScratch indicates that a vm uses the from_scratch template type.
	TemplateTypeFromScratch TemplateType = "from_scratch"
	// TemplateTypeTemplates indicates that a vm uses the templates template type.
	TemplateTypeTemplates TemplateType = "templates"
)

// ProvisioningRequest represents a VM provisioning request.
type ProvisioningRequest struct {
	// required fields
	Hostname string `json:"hostname"`

	// optional fields
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
	VideoMemoryMB      *int                                `json:"video_memory_mb,omitempty"`
	DNS1               *string                             `json:"dns1,omitempty"`
	DNS2               *string                             `json:"dns2,omitempty"`
	DNS3               *string                             `json:"dns3,omitempty"`
	DNS4               *string                             `json:"dns4,omitempty"`
	Password           *string                             `json:"password,omitempty"`
	SSH                *string                             `json:"ssh,omitempty"`
	Script             *string                             `json:"script,omitempty"`
	BootDelay          *int                                `json:"boot_delay,omitempty"`
	EnterBiosSetup     *bool                               `json:"enter_bios_setup,omitempty"`
	Organization       *string                             `json:"organization,omitempty"`
	CustomName         *string                             `json:"custom_name,omitempty"`
	VTPMEnabled        *bool                               `json:"vtpm_enabled,omitempty"`
	Firmware           *string                             `json:"firmware,omitempty"`
	OSHostname         *string                             `json:"os_hostname,omitempty"`
}

// ProvisioningResponse represents the response of a provisioning request.
type ProvisioningResponse struct {
	Progress       int      `json:"progress"`
	Errors         []string `json:"errors"`
	TaskIdentifier string   `json:"identifier"`
	Queued         bool     `json:"queued"`
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
	VLan           string   `json:"vlan"`
	IPs            []string `json:"ips"`
}

// TemplateResponse represents the templates from a location.
type TemplateResponse struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Architecture string         `json:"architecture"`
	Bit          string         `json:"bit"`
	Build        string         `json:"build"`
	Params       map[string]any `json:"params"`
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
	err := c.transport.GetSingle(ctx, "/api/vsphere/v1/provisioning/cpu_architecture.json", &resp)
	return resp, common.MapTransportError(err)
}

// GetCPUPerformanceTypes returns all available cpu performance types.
func (c *ProvisioningClient) GetCPUPerformanceTypes(ctx context.Context) ([]CPUPerformanceType, error) {
	var resp []CPUPerformanceType
	err := c.transport.GetSingle(ctx, "/api/vsphere/v1/provisioning/cpu_performance_type.json", &resp)
	return resp, common.MapTransportError(err)
}

// GetDiskTypes returns all available disk types in a location.
func (c *ProvisioningClient) GetDiskTypes(ctx context.Context, locationIdentifier string) ([]DiskType, error) {
	var resp []DiskType
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/vsphere/v1/provisioning/disk_type.json/%s", locationIdentifier), &resp)
	return resp, common.MapTransportError(err)
}

// ListLocations lists a paged response of locations.
func (c *ProvisioningClient) ListLocations(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[Location], error) {
	resp := paging.PagedResponse[Location]{}
	err := c.transport.Get(ctx, "/api/vsphere/v1/provisioning/location.json", &resp, pageParams, nil)
	return resp, common.MapTransportError(err)
}

// ListLocationPageFetcher returns a paging.PageFetcher for locations.
func (c *ProvisioningClient) ListLocationPageFetcher() paging.PageFetcher[Location] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[Location], error) {
		return c.ListLocations(ctx, pageParams)
	}
}

// ListTemplates returns a paging.PageFetcher for templates.
func (c *ProvisioningClient) ListTemplates(ctx context.Context,
	locationIdentifier string, templateType TemplateType) ([]TemplateResponse, error) {
	var resp []TemplateResponse
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/vsphere/v1/provisioning/templates.json/%s/%s", locationIdentifier, templateType), &resp)
	return resp, common.MapTransportError(err)
}

// ListAvailabilityZones lists a paged response of availability zones in a location.
func (c *ProvisioningClient) ListAvailabilityZones(ctx context.Context, locationIdentifier string) ([]AvailabilityZone, error) {
	var resp []AvailabilityZone
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/vsphere/v1/provisioning/location.json/%s/availability_zone", locationIdentifier), &resp)
	return resp, common.MapTransportError(err)
}

// GetNicTypes returns all available nic types.
func (c *ProvisioningClient) GetNicTypes(ctx context.Context) ([]NicType, error) {
	var resp []NicType
	err := c.transport.GetSingle(ctx, "/api/vsphere/v1/provisioning/nic_type.json", &resp)
	return resp, common.MapTransportError(err)
}

// GetProvisioningProgress returns the progress for the specified vm provisioning.
func (c *ProvisioningClient) GetProvisioningProgress(ctx context.Context, taskIdentifier string) (ProvisioningProgress, error) {
	resp := ProvisioningProgress{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/vsphere/v1/provisioning/progress.json/%s", taskIdentifier), &resp)
	return resp, common.MapTransportError(err)
}

// AwaitCompletion polls the status of a started provisioning request and blocks until it is done.
//
// ctx will be checked for cancellation and the method returns immediately if so.
// taskIdentifier is the running provisioning task and is contained within ProvisioningResponse.
//
// Returned will be the VM identifier and an error if polling or provision failed.
func (c *ProvisioningClient) AwaitCompletion(ctx context.Context, taskIdentifier string) (string, error) {
	const (
		pollInterval          = 10 * time.Second
		progressCompleteValue = 100
	)

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			progressResponse, err := c.GetProvisioningProgress(ctx, taskIdentifier)
			switch {
			case common.IsNotFoundError(err):
				return "", fmt.Errorf("could not get provision progress. not found: %w", err)
			case err == nil:
				if progressResponse.Progress == progressCompleteValue {
					return progressResponse.VMIdentifier, nil
				}
			default:
				return "", fmt.Errorf("could not get provision progress: %w", err)
			}
		case <-ctx.Done():
			return "", fmt.Errorf("vm did not get ready in time: %w", ctx.Err())
		}
	}
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
	err := c.transport.Post(ctx, fmt.Sprintf("/api/vsphere/v1/provisioning/vm.json/%s/%s/%s", locationIdentifier, templateType, templateIdentifier), request, &resp)
	return resp, common.MapTransportError(err)
}

// ProvisionFromScratch provisions a new vm from scratch.
func (c *ProvisioningClient) ProvisionFromScratch(
	ctx context.Context,
	locationIdentifier string,
	templateIdentifier string,
	request ProvisioningRequest,
) (ProvisioningResponse, error) {
	return c.Provision(ctx, locationIdentifier, TemplateTypeFromScratch, templateIdentifier, request)
}

// ProvisionTemplate provisions a new vm using a template.
func (c *ProvisioningClient) ProvisionTemplate(
	ctx context.Context,
	locationIdentifier string,
	templateIdentifier string,
	request ProvisioningRequest,
) (ProvisioningResponse, error) {
	return c.Provision(ctx, locationIdentifier, TemplateTypeTemplates, templateIdentifier, request)
}
