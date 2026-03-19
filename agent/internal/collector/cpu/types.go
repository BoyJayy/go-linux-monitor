package cpu

type CPUCounters struct {
	User    uint64
	Nice    uint64
	System  uint64
	Idle    uint64
	Iowait  uint64
	Irq     uint64
	Softirq uint64
	Steal   uint64
	Total   uint64
}

type Usage struct {
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
	PerCorePct map[string]Usage
}
