package snapshot

import (
	"monitoring/api"
	"time"
)

type CPUResult struct {
	stat api.UsageWithCores
	err  error
}

type MemResult struct {
	stat api.MemoryUsage
	err  error
}

type DiskResult struct {
	stat []api.DiskUsage
	err  error
}
type NetworkResult struct {
	stat api.NetworkUsage
	err  error
}

type Metrics struct {
	Timestamp time.Time          `json:"timestamp"`
	HostID    string             `json:"host_id"`
	CPU       api.UsageWithCores `json:"cpu"`
	Mem       api.MemoryUsage    `json:"mem"`
	Disk      []api.DiskUsage    `json:"disk"`
	Network   api.NetworkUsage   `json:"network"`
}
