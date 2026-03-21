package memory

type MemCounters struct {
	Total     uint64
	Available uint64
	Used      uint64
}

type MemoryUsage struct {
	TotalBytes     uint64
	AvailableBytes uint64
	UsedBytes      uint64
	UsedPct        float64
}
