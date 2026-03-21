package memory

func Collect() (MemoryUsage, error) {
	cur, err := ReadMemCounters()
	if err != nil {
		return MemoryUsage{}, err
	}
	return MemoryUsage{
		TotalBytes:     cur.Total,
		AvailableBytes: cur.Available,
		UsedBytes:      cur.Used,
		UsedPct:        float64(cur.Used) / float64(cur.Total) * 100,
	}, nil
}
