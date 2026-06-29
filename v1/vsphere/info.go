package vsphere

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/v1/common"
)

// InfoGetResponse represents the response of the info get endpoint.
type InfoGetResponse struct {
	Identifier                     string                          `json:"identifier"`
	Name                           string                          `json:"name"`
	CustomName                     string                          `json:"custom_name"`
	GuestOS                        string                          `json:"guest_os"`
	Firmware                       string                          `json:"firmware"`
	RAM                            int                             `json:"ram"`
	CPU                            int                             `json:"cpu"`
	CPUClockRate                   int                             `json:"cpu_clock_rate"`
	CPUPerformanceType             string                          `json:"cpu_performance_type"`
	VTPMEnabled                    bool                            `json:"vtpm_enabled"`
	Cores                          int                             `json:"cores"`
	Disks                          int                             `json:"disks"`
	DiskInfo                       []InfoGetResponseDiskInfo       `json:"disk_info"`
	Network                        []InfoGetResponseNetwork        `json:"network"`
	CDRom                          string                          `json:"cdrom"`
	VersionTools                   string                          `json:"version_tools"`
	GuestToolsStatus               string                          `json:"guest_tools_status"`
	LocationCode                   string                          `json:"location_code"`
	LocationCountry                string                          `json:"location_country"`
	LocationIdentifier             string                          `json:"location_identifier"`
	LocationName                   string                          `json:"location_name"`
	ProvisioningLocationIdentifier string                          `json:"provisioning_location_identifier"`
	TemplateID                     string                          `json:"template_id"`
	ResourceSalesperson            string                          `json:"resource_salesperson"`
	AvailabilityZone               InfoGetResponseAvailabilityZone `json:"availability_zone"`
}

// InfoGetResponseDiskInfo represents infos about a single disk in the get response.
type InfoGetResponseDiskInfo struct {
	DiskGb       int    `json:"disk_gb"`
	DiskID       int    `json:"disk_id"`
	DiskType     string `json:"disk_type"`
	Iops         int    `json:"iops"`
	Latence      int    `json:"latence"`
	StorageType  string `json:"storage_type"`
	BusType      string `json:"bus_type"`
	BusTypeLabel string `json:"bus_type_label"`
}

// InfoGetResponseNetwork represents infos about a single network in the get response.
type InfoGetResponseNetwork struct {
	Nic            int      `json:"nic"`
	NicType        string   `json:"nic_type"`
	BandwidthLimit int      `json:"bandwidth_limit"`
	Vlan           string   `json:"vlan"`
	ID             int      `json:"id"`
	IpsV4          []string `json:"ips_v4"`
	IpsV6          []string `json:"ips_v6"`
	MacAddress     string   `json:"mac_address"`
	Mode           string   `json:"mode"`
	Vlans          []int    `json:"vlans"`
}

// InfoGetResponseAvailabilityZone represents the info about the availability zone in the get response.
type InfoGetResponseAvailabilityZone struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// InfoClient is an api client for vm infos.
type InfoClient struct {
	transport *internal.Transport
}

func newInfoClient(transport *internal.Transport) *InfoClient {
	return &InfoClient{
		transport: transport,
	}
}

// Get returns an info object by identifier.
func (c *InfoClient) Get(ctx context.Context, identifier string) (*InfoGetResponse, error) {
	resp := InfoGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/vsphere/v1/info.json/%s/info", identifier), &resp)
	return &resp, common.MapTransportError(err)
}
