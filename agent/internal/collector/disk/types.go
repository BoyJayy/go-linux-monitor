package disk

type DiskCounters struct {
	Total uint64
	Free  uint64
	Used  uint64
}

type DiskUsage struct {
	Mount      string  
	TotalBytes uint64  
	FreeBytes  uint64  
	UsedBytes  uint64  
	UsedPct    float64 
}