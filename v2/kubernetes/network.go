package kubernetes

import (
	"time"

	"github.com/anexia/go-anxsdk/v2/common"
)

type NetworkBandwidthLimit string

const (
	NetworkBandwidthLimit100Mbit = NetworkBandwidthLimit("100")
	NetworkBandwidthLimit1Gbit   = NetworkBandwidthLimit("1000")
	NetworkBandwidthLimit10Gbit  = NetworkBandwidthLimit("10000")
)

// NetworkGetResponse represents the networks of a Nodepool.
type NetworkGetResponse struct {
	Identifier string          `json:"identifier,omitempty"`
	Name       string          `json:"name,omitempty"`
	Nodepool   common.Resource `json:"nodepool,omitempty"`

	BandwidthLimit common.IDTitleTuple[NetworkBandwidthLimit] `json:"bandwidth_limit,omitempty"`
	VLAN           common.Resource                            `json:"vlan,omitempty"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// NetworkUpdateRequest represents the networks of a Nodepool.
type NetworkUpdateRequest struct {
	Name       string `json:"name,omitempty"`
	NodepoolID string `json:"nodepool,omitempty"`

	BandwidthLimit NetworkBandwidthLimit `json:"bandwidth_limit,omitempty"`
	VLANID         string                `json:"vlan,omitempty"`
}
