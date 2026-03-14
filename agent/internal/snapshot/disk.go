package snapshot

import "golang.org/x/sys/unix"

func ReadDiskSnapshot() (snap DiskSnapshot, err error) {
	var stat unix.Statfs_t
	err = unix.Statfs("/", &stat)
	if err != nil {
		return DiskSnapshot{}, err
	}
	total := uint64(stat.Blocks) * uint64(stat.Bsize)
	free := uint64(stat.Bavail) * uint64(stat.Bsize)
	used := total - free
	return DiskSnapshot{
		Total: total,
		Free:  free,
		Used:  used,
	}, nil

}
