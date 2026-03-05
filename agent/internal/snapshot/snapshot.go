package snapshot

import "time"

type Snapshot struct {
	Device_id         string
	Timestamp         time.Time
	Cpu_usage_percent float64
	Mem_total_bytes   uint64
	Mem_used_bytes    uint64
}

/*type Snapshot1 struct {
}*/
