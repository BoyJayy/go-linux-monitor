package snapshot

import "time"

type CPUStat struct {
	TotalPct   float64 `json:"total_pct"`
	UserPct    float64 `json:"user_pct"`
	SystemPct  float64 `json:"system_pct"`
	NicePct    float64 `json:"nice_pct"`
	IowaitPct  float64 `json:"iowait_pct"`
	IrqPct     float64 `json:"irq_pct"`
	SoftirqPct float64 `json:"softirq_pct"`
	StealPct   float64 `json:"steal_pct"`
}

type MemStat struct {
	TotalBytes     uint64  `json:"total_bytes"`
	AvailableBytes uint64  `json:"available_bytes"`
	UsedBytes      uint64  `json:"used_bytes"`
	UsedPct        float64 `json:"used_pct"`
}

type DiskStat struct {
	Mount      string  `json:"mount"`
	TotalBytes uint64  `json:"total_bytes"`
	FreeBytes  uint64  `json:"free_bytes"`
	UsedBytes  uint64  `json:"used_bytes"`
	UsedPct    float64 `json:"used_pct"`
}

type NetIfaceStat struct {
	Name    string  `json:"name"`
	RxBytes uint64  `json:"rx_bytes"`
	TxBytes uint64  `json:"tx_bytes"`
	RxBps   float64 `json:"rx_bps"`
	TxBps   float64 `json:"tx_bps"`
}

type NetworkStat struct {
	RxBytesTotal uint64         `json:"rx_bytes_total"`
	TxBytesTotal uint64         `json:"tx_bytes_total"`
	RxBpsTotal   float64        `json:"rx_bps_total"`
	TxBpsTotal   float64        `json:"tx_bps_total"`
	Ifaces       []NetIfaceStat `json:"ifaces,omitempty"`
}

type Metrics struct {
	Timestamp time.Time   `json:"timestamp"`
	HostID    string      `json:"host_id"`
	CPU       CPUStat     `json:"cpu"`
	Mem       MemStat     `json:"mem"`
	Disk      []DiskStat  `json:"disk"`
	Network   NetworkStat `json:"network"`
}
