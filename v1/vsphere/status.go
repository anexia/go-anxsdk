package vsphere

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/v1/common"
)

// StatusGetResponse represents the response of the status get endpoint.
type StatusGetResponse struct {
	Identifier     string                        `json:"identifier"`
	CPUUsage       int                           `json:"cpu_usage"`
	MemoryUsage    int                           `json:"memory_usage"`
	DiskSpaceTotal int                           `json:"disk_space_total"`
	DiskSpaceFree  int                           `json:"disk_space_free"`
	MountPoints    []StatusGetResponseMountPoint `json:"mount_points"`
	LastRefresh    StatusGetResponseLastRefresh  `json:"last_refresh"`
	Total          string                        `json:"total"`
	Free           string                        `json:"free"`
	Used           string                        `json:"used"`
	Percentage     float32                       `json:"percentage"`
	DiskSizeState  string                        `json:"disk_size_state"`
	PowerState     string                        `json:"power_state"`
}

// StatusGetResponseMountPoint represents infos about a single mount point.
type StatusGetResponseMountPoint struct {
	Name             string  `json:"name"`
	CapacityInBytes  string  `json:"capacity_in_bytes"`
	FreeSpaceInBytes string  `json:"free_space_in_bytes"`
	UsedSpace        int     `json:"used_space"`
	Percentage       float64 `json:"percentage"`
	Total            string  `json:"total"`
	Free             string  `json:"free"`
	Used             string  `json:"used"`
	State            string  `json:"state"`
}

// StatusGetResponseLastRefresh represents infos about the last refresh.
type StatusGetResponseLastRefresh struct {
	Date string `json:"date"`
	Day  string `json:"day"`
	Time string `json:"time"`
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
