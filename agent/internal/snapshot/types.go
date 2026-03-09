package snapshot

import (
	"time"
)

type MemSnapshot struct {
	Total     uint64
	Available uint64
	Used      uint64
}

type CpuSnapshot struct {
	User    uint64 `json:"mount"`
	Nice    uint64 `json:"nice"`
	System  uint64 `json:"system"`
	Idle    uint64 `json:"idle"`
	Iowait  uint64 `json:"iowait"`
	Irq     uint64 `json:"irq"`
	Softirq uint64 `json:"softirq"`
	Steal   uint64 `json:"steal"`
	Total   uint64 `json:"total"`
}

type Snapshot struct {
	DeviceID        string
	Timestamp       time.Time
	CpuUsagePercent float64
	MemTotaBytes    uint64
	MemUsedBytes    uint64
}
