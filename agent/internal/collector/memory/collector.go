package memory

import "monitoring/api"

func Collect() (api.MemoryUsage, error) {
	cur, err := ReadMemCounters()
	if err != nil {
		return api.MemoryUsage{}, err
	}
	return api.MemoryUsage{
		TotalBytes:     cur.Total,
		AvailableBytes: cur.Available,
		UsedBytes:      cur.Used,
		UsedPct:        float64(cur.Used) / float64(cur.Total) * 100,
	}, nil
}
