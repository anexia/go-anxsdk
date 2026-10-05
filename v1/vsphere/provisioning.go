package vsphere

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

var (
	// ErrNegativeValue indicates a numeric field was given a negative value.
	ErrNegativeValue = errors.New("must not be negative")
	// ErrMissingCredentials indicates that neither an ssh key nor a password was provided.
	ErrMissingCredentials = errors.New("must provide either SSH or Password")
	// ErrInvalidIP indicates that a field could not be parsed as an IP address.
	ErrInvalidIP = errors.New("not a valid IP")
)

// CPUArchitecture represents the cpu architecture of a vm.
type CPUArchitecture struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// CPUPerformanceTypeResponse represents the cpu performance characteristics.
type CPUPerformanceTypeResponse struct {
	ID             string  `json:"id"`
	Architecture   *string `json:"architecture,omitempty"`
	Prioritization string  `json:"prioritization"`
	Limit          float32 `json:"limit"`
	Unit           string  `json:"unit"`
}

// DiskTypeResponse represents the disk characteristics.
type DiskTypeResponse struct {
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

// Firmware is the firmware interface a vm boots with.
type Firmware string

const (
	// FirmwareBIOS indicates that a vm boots using legacy BIOS.
	FirmwareBIOS Firmware = "BIOS"
	// FirmwareUEFI indicates that a vm boots using UEFI.
	FirmwareUEFI Firmware = "UEFI"
)

// CPUPerformanceType is the cpu performance class assigned to a vm.
type CPUPerformanceType string

// The available cpu performance types, per cpu vendor.
const (
	CPUPerformanceTypeBestEffort      CPUPerformanceType = "best-effort"
	CPUPerformanceTypeStandard        CPUPerformanceType = "standard"
	CPUPerformanceTypeEnterprise      CPUPerformanceType = "enterprise"
	CPUPerformanceTypePerformance     CPUPerformanceType = "performance"
	CPUPerformanceTypePerformancePlus CPUPerformanceType = "performance-plus"

	CPUPerformanceTypeBestEffortIntel      CPUPerformanceType = "best-effort-intel"
	CPUPerformanceTypeStandardIntel        CPUPerformanceType = "standard-intel"
	CPUPerformanceTypeEnterpriseIntel      CPUPerformanceType = "enterprise-intel"
	CPUPerformanceTypePerformanceIntel     CPUPerformanceType = "performance-intel"
	CPUPerformanceTypePerformancePlusIntel CPUPerformanceType = "performance-plus-intel"

	CPUPerformanceTypeBestEffortAMD  CPUPerformanceType = "best-effort-amd"
	CPUPerformanceTypeStandardAMD    CPUPerformanceType = "standard-amd"
	CPUPerformanceTypeEnterpriseAMD  CPUPerformanceType = "enterprise-amd"
	CPUPerformanceTypePerformanceAMD CPUPerformanceType = "performance-amd"
)

// DiskType is the performance class of a vm disk.
type DiskType string

// The available disk types: enterprise (ENT), high performance compute (HPC), local and standard (STD).
const (
	DiskTypeENT1 DiskType = "ENT1"
	DiskTypeENT2 DiskType = "ENT2"
	DiskTypeENT3 DiskType = "ENT3"
	DiskTypeENT4 DiskType = "ENT4"
	DiskTypeENT5 DiskType = "ENT5"
	DiskTypeENT6 DiskType = "ENT6"

	DiskTypeHPC1 DiskType = "HPC1"
	DiskTypeHPC2 DiskType = "HPC2"
	DiskTypeHPC3 DiskType = "HPC3"
	DiskTypeHPC4 DiskType = "HPC4"
	DiskTypeHPC5 DiskType = "HPC5"

	DiskTypeLocal DiskType = "LOC3"

	DiskTypeSTD1 DiskType = "STD1"
	DiskTypeSTD2 DiskType = "STD2"
	DiskTypeSTD3 DiskType = "STD3"
	DiskTypeSTD4 DiskType = "STD4"
	DiskTypeSTD5 DiskType = "STD5"
	DiskTypeSTD6 DiskType = "STD6"
)

// ProvisioningRequest represents a VM provisioning request.
type ProvisioningRequest struct {
	// required fields
	Hostname string `json:"hostname"`

	// optional fields
	MemoryMB           *int                                `json:"memory_mb,omitempty"`
	CPUs               *int                                `json:"cpus,omitempty"`
	DiskGB             *int                                `json:"disk_gb,omitempty"`
	DiskType           *DiskType                           `json:"disk_type,omitempty"`
	AdditionalDisks    []ProvisioningRequestAdditionalDisk `json:"additional_disks,omitempty"`
	CPUPerformanceType *CPUPerformanceType                 `json:"cpu_performance_type,omitempty"`
	AvailabilityZone   *string                             `json:"availability_zone,omitempty"`
	Sockets            *int                                `json:"sockets,omitempty"`
	Network            []ProvisioningRequestNetwork        `json:"network,omitempty"`
	VideoMemoryAuto    *bool                               `json:"video_memory_auto,omitempty"`
	VideoMemoryMB      *int                                `json:"video_memory_mb,omitempty"`
	DNS1               *string                             `json:"dns1,omitempty"`
	DNS2               *string                             `json:"dns2,omitempty"`
	DNS3               *string                             `json:"dns3,omitempty"`
	DNS4               *string                             `json:"dns4,omitempty"`
	Password           *string                             `json:"password,omitempty"`
	SSH                *string                             `json:"ssh,omitempty"`
	Script             *string                             `json:"script,omitempty"`
	BootDelaySeconds   *int                                `json:"boot_delay,omitempty"`
	EnterBiosSetup     *bool                               `json:"enter_bios_setup,omitempty"`
	Organization       *string                             `json:"organization,omitempty"`
	CustomName         *string                             `json:"custom_name,omitempty"`
	VTPMEnabled        *bool                               `json:"vtpm_enabled,omitempty"`
	Firmware           *Firmware                           `json:"firmware,omitempty"`
	OSHostname         *string                             `json:"os_hostname,omitempty"`
}

// validateNonNegative returns an error if value is set and negative.
func validateNonNegative(field string, value *int) error {
	if value != nil && *value < 0 {
		return fmt.Errorf("%s %w: %d", field, ErrNegativeValue, *value)
	}

	return nil
}

// validateOptionalIP returns an error if value is set to a non-empty string that is not an IP address.
func validateOptionalIP(field string, value *string) error {
	if value != nil && *value != "" && net.ParseIP(*value) == nil {
		return fmt.Errorf("%s is %w: %q", field, ErrInvalidIP, *value)
	}

	return nil
}

// PreValidate checks the request for obvious problems before it is sent to the api.
// Unset optional fields are not validated. All violations are reported together via errors.Join.
func (r *ProvisioningRequest) PreValidate() error {
	var errs []error

	if r.Script != nil && *r.Script != "" {
		if _, err := base64.StdEncoding.DecodeString(*r.Script); err != nil {
			errs = append(errs, fmt.Errorf("Script is not base64 encoded: %w", err)) //nolint:staticcheck
		}
	}

	hasSSH := r.SSH != nil && *r.SSH != ""
	hasPassword := r.Password != nil && *r.Password != ""
	if !hasSSH && !hasPassword {
		errs = append(errs, ErrMissingCredentials)
	}

	errs = append(errs,
		validateNonNegative("MemoryMB", r.MemoryMB),
		validateNonNegative("CPUs", r.CPUs),
		validateNonNegative("DiskGB", r.DiskGB),
		validateNonNegative("Sockets", r.Sockets),
		validateNonNegative("BootDelaySeconds", r.BootDelaySeconds),
		validateOptionalIP("DNS1", r.DNS1),
		validateOptionalIP("DNS2", r.DNS2),
		validateOptionalIP("DNS3", r.DNS3),
		validateOptionalIP("DNS4", r.DNS4),
	)

	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("provisioning request validation: %w", err)
	}

	return nil
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

// BuildNumber returns the parsed build number as int.
func (t *TemplateResponse) BuildNumber() (int, error) {
	if t.Build == "" || t.Build[0] != 'b' {
		return 0, errors.New("template build does not start with 'b'") //nolint:err113,wrapcheck
	}

	buildNumber, err := strconv.Atoi(t.Build[1:])
	if err != nil {
		return 0, fmt.Errorf("template build cannot be parsed: %w", err)
	}

	return buildNumber, nil
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
func (c *ProvisioningClient) GetCPUPerformanceTypes(ctx context.Context) ([]CPUPerformanceTypeResponse, error) {
	var resp []CPUPerformanceTypeResponse
	err := c.transport.GetSingle(ctx, "/api/vsphere/v1/provisioning/cpu_performance_type.json", &resp)
	return resp, common.MapTransportError(err)
}

// GetDiskTypes returns all available disk types in a location.
func (c *ProvisioningClient) GetDiskTypes(ctx context.Context, locationIdentifier string) ([]DiskTypeResponse, error) {
	var resp []DiskTypeResponse
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
func (c *ProvisioningClient) ListTemplates(ctx context.Context, locationIdentifier string, templateType TemplateType) ([]TemplateResponse, error) {
	var resp []TemplateResponse
	err := c.transport.Get(ctx, fmt.Sprintf("/api/vsphere/v1/provisioning/templates.json/%s/%s", locationIdentifier, templateType), &resp, paging.NewParams(1, 1000), nil) //nolint:revive
	return resp, common.MapTransportError(err)
}

const (
	// LatestTemplateBuild is used to find the template with the highest build number.
	LatestTemplateBuild = "latest"
)

// FindNamedTemplate retrieves a template by name and build at a specified location.
// Empty and LatestTemplateBuild build identifier will yield the highest available build.
func (c *ProvisioningClient) FindNamedTemplate(ctx context.Context, locationIdentifier, name, build string) (*TemplateResponse, error) { //nolint:revive // it's not that hard
	var match *TemplateResponse
	buildNo := -1
	useLatest := build == "" || build == LatestTemplateBuild

	allTemplates, err := c.ListTemplates(ctx, locationIdentifier, TemplateTypeTemplates)
	if err != nil {
		return nil, fmt.Errorf("error listing templates: %w", err)
	}

	for _, tmpl := range allTemplates {
		if tmpl.Name != name {
			continue
		}

		if useLatest {
			currentTemplateBuildNo, err := tmpl.BuildNumber()
			if err != nil {
				continue
			}

			if match == nil || currentTemplateBuildNo > buildNo {
				match = &tmpl
				buildNo = currentTemplateBuildNo
			}
		} else if tmpl.Build == build {
			match = &tmpl
			break
		}
	}

	if match == nil {
		return nil, fmt.Errorf("%w: named template not found: name=%q, build=%q, location=%q", &common.APIError{StatusCode: http.StatusNotFound}, name, build, locationIdentifier)
	}

	return match, nil
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

// DeprovisioningResponse represents the response of a vm deprovisioning request.
type DeprovisioningResponse struct {
	Identifier             string `json:"identifier"`
	DeleteWillBeExecutedAt string `json:"delete_will_be_executed_at"`
}

// deprovisionParams defines the available query parameters for the vm deprovisioning endpoint.
type deprovisionParams struct {
	Delayed bool `url:"delayed"`
}

// Deprovision issues a request to deprovision an existing vm.
//
// identifier is the vm identifier, as returned in ProvisioningProgress.VMIdentifier once provisioning completed.
// delayed indicates that the vm shall be removed with a delay of 24h instead of immediately.
func (c *ProvisioningClient) Deprovision(ctx context.Context, identifier string, delayed bool) (DeprovisioningResponse, error) {
	resp := DeprovisioningResponse{}
	endpoint := fmt.Sprintf("/api/vsphere/v1/provisioning/vm.json/%s", identifier)
	err := c.transport.DeleteWithResponse(ctx, endpoint, deprovisionParams{Delayed: delayed}, &resp)
	return resp, common.MapTransportError(err)
}
