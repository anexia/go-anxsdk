package kubernetes

import (
	"time"

	"github.com/anexia/go-anxsdk/v2/common"
)

type DiskPerformanceType string

const (
	DiskPerformanceTypeSTD1 DiskPerformanceType = "STD1"
	DiskPerformanceTypeSTD2 DiskPerformanceType = "STD2"
	DiskPerformanceTypeSTD3 DiskPerformanceType = "STD3"
	DiskPerformanceTypeSTD4 DiskPerformanceType = "STD4"
	DiskPerformanceTypeSTD5 DiskPerformanceType = "STD5"
	DiskPerformanceTypeENT1 DiskPerformanceType = "ENT1"
	DiskPerformanceTypeENT2 DiskPerformanceType = "ENT2"
	DiskPerformanceTypeENT3 DiskPerformanceType = "ENT3"
	DiskPerformanceTypeENT4 DiskPerformanceType = "ENT4"
	DiskPerformanceTypeENT5 DiskPerformanceType = "ENT5"
	DiskPerformanceTypeENT6 DiskPerformanceType = "ENT6"

	DiskPerformanceTypeDefault = DiskPerformanceTypeENT6
)

// DisksGetResponse represents the disks of a Nodepool.
type DisksGetResponse struct {
	Identifier string          `json:"identifier,omitempty"`
	Name       string          `json:"name,omitempty"`
	Nodepool   common.Resource `json:"nodepool,omitempty"`

	SizeBytes       uint64                                   `json:"size_bytes,omitempty"`
	PerformanceType common.IDTitleTuple[DiskPerformanceType] `json:"performance_type,omitempty"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// DiskUpdateRequest represents the disks of a Nodepool.
type DiskUpdateRequest struct {
	Name       string `json:"name,omitempty"`
	NodepoolID string `json:"nodepool,omitempty"`

	SizeBytes       uint64              `json:"size_bytes,omitempty"`
	PerformanceType DiskPerformanceType `json:"performance_type,omitempty"`
}
