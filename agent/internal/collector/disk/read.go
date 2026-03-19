package disk

import "golang.org/x/sys/unix"

func ReadDiskSnapshot() (snap DiskCounters, err error) {
	var stat unix.Statfs_t
	err = unix.Statfs("/", &stat)
	if err != nil {
		return DiskCounters{}, err
	}
	total := uint64(stat.Blocks) * uint64(stat.Bsize)
	free := uint64(stat.Bavail) * uint64(stat.Bsize)
	used := total - free
	return DiskCounters{
		Total: total,
		Free:  free,
		Used:  used,
	}, nil
}
