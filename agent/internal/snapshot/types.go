package snapshot

import (
	"agent/internal/collector/cpu"
	"agent/internal/collector/disk"
	"agent/internal/collector/memory"
	"agent/internal/collector/network"
	"time"
)

type CPUResult struct {
	stat cpu.UsageWithCores
	err  error
}

type MemResult struct {
	stat memory.MemoryUsage
	err  error
}

type DiskResult struct {
	stat []disk.DiskUsage
	err  error
}
type NetworkResult struct {
	stat network.NetworkUsage
	err  error
}

type Metrics struct {
	Timestamp time.Time            `json:"timestamp"`
	HostID    string               `json:"host_id"`
	CPU       cpu.UsageWithCores   `json:"cpu"`
	Mem       memory.MemoryUsage   `json:"mem"`
	Disk      []disk.DiskUsage     `json:"disk"`
	Network   network.NetworkUsage `json:"network"`
}
