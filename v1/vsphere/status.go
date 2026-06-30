package vsphere

import (
	"context"
	"fmt"
	"time"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/v1/common"
)

// PowerState represents the possible power states.
type PowerState string

const (
	// PowerStatePoweredOn shows that a VM is running.
	PowerStatePoweredOn PowerState = "PoweredOn"
	// PowerStatePoweredOff shows that a VM is turned off.
	PowerStatePoweredOff PowerState = "PoweredOff"
)

// StatusGetResponse represents the response of the status get endpoint.
type StatusGetResponse struct {
	Identifier          string                        `json:"identifier"`
	CPUUsagePercent     float32                       `json:"cpu_usage_percent"`
	MemoryUsageMegabyte int                           `json:"memory_usage"`
	DiskSpaceTotalBytes string                        `json:"disk_space_total"`
	DiskSpaceFreeBytes  string                        `json:"disk_space_free"`
	MountPoints         []StatusGetResponseMountPoint `json:"mount_points"`
	LastRefresh         StatusGetResponseLastRefresh  `json:"last_refresh"`
	PowerState          PowerState                    `json:"power_state"`
}

// StatusGetResponseMountPoint represents infos about a single mount point.
type StatusGetResponseMountPoint struct {
	// Name is the mountpoint directory.
	Name             string  `json:"name"`
	CapacityInBytes  uint64  `json:"capacity_in_bytes"`
	FreeSpaceInBytes uint64  `json:"free_space_in_bytes"`
	UsedSpace        uint64  `json:"used_space"`
	Percentage       float64 `json:"percentage"`
}

// StatusGetResponseLastRefresh represents infos about the last refresh.
type StatusGetResponseLastRefresh struct {
	Date time.Time `json:"date"`
	Day  string    `json:"day"`
	Time string    `json:"time"`
}

// StatusClient is an api client for getting vm status.
type StatusClient struct {
	transport *internal.Transport
}

func newStatusClient(transport *internal.Transport) *StatusClient {
	return &StatusClient{
		transport: transport,
	}
}

// Get returns a status object by identifier.
func (c *StatusClient) Get(ctx context.Context, identifier string) (StatusGetResponse, error) {
	resp := StatusGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/vshpere/v1/status.json/%s/info", identifier), &resp)
	return resp, common.MapTransportError(err)
}
