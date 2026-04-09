package api

import "time"

type MemoryUsage struct {
	TotalBytes     uint64
	AvailableBytes uint64
	UsedBytes      uint64
	UsedPct        float64
}

type CPUUsage struct {
	UserPct    float64
	NicePct    float64
	SystemPct  float64
	IdlePct    float64
	IowaitPct  float64
	IrqPct     float64
	SoftirqPct float64
	StealPct   float64
	TotalPct   float64
}

type DiskUsage struct {
	Mount      string
	TotalBytes uint64
	FreeBytes  uint64
	UsedBytes  uint64
	UsedPct    float64
}

type UsageWithCores struct {
	UserPct    float64
	NicePct    float64
	SystemPct  float64
	IdlePct    float64
	IowaitPct  float64
	IrqPct     float64
	SoftirqPct float64
	StealPct   float64
	TotalPct   float64
	PerCorePct map[string]CPUUsage
}

type NetIfaceStat struct {
	Name    string
	RxBytes uint64
	TxBytes uint64
	RxBps   float64
	TxBps   float64
}

type NetworkUsage struct {
	RxBytesTotal uint64
	TxBytesTotal uint64
	RxBpsTotal   float64
	TxBpsTotal   float64
	Ifaces       []NetIfaceStat
}

type CPUResult struct {
	stat UsageWithCores
	err  error
}

type MemResult struct {
	stat MemoryUsage
	err  error
}

type DiskResult struct {
	stat []DiskUsage
	err  error
}
type NetworkResult struct {
	stat NetworkUsage
	err  error
}

type Metrics struct {
	Timestamp time.Time      `json:"timestamp"`
	HostID    string         `json:"host_id"`
	CPU       UsageWithCores `json:"cpu"`
	Mem       MemoryUsage    `json:"mem"`
	Disk      []DiskUsage    `json:"disk"`
	Network   NetworkUsage   `json:"network"`
}
