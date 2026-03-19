package cpu

type CPUCounters struct {
	User    uint64 `json:"user"`
	Nice    uint64 `json:"nice"`
	System  uint64 `json:"system"`
	Idle    uint64 `json:"idle"`
	Iowait  uint64 `json:"iowait"`
	Irq     uint64 `json:"irq"`
	Softirq uint64 `json:"softirq"`
	Steal   uint64 `json:"steal"`
	Total   uint64 `json:"total"`
}
