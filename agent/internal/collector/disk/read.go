package disk

import "golang.org/x/sys/unix"

func ReadDiskCounters(mountpoint string) (snap DiskCounters, err error) {
	var stat unix.Statfs_t
	err = unix.Statfs(mountpoint, &stat)
	if err != nil {
		return DiskCounters{}, err
	}
	total := uint64(stat.Blocks) * uint64(stat.Bsize)
	free := uint64(stat.Bfree) * uint64(stat.Bsize)
	used := total - free
	return DiskCounters{
		Total: total,
		Free:  free,
		Used:  used,
	}, nil
}

func ConvertFromCountersToUsage(mountpoint string, snap DiskCounters) DiskUsage {
	pct := float64(0)
	if snap.Total != 0 {
		pct = float64(snap.Used) / float64(snap.Total) * 100
	}
	return DiskUsage{
		Mount:      mountpoint,
		TotalBytes: snap.Total,
		FreeBytes:  snap.Free,
		UsedBytes:  snap.Used,
		UsedPct:    pct,
	}
}
