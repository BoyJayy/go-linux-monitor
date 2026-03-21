package disk

import (
	"errors"
	"os"
	"strings"
)

var skippedTypes = map[string]struct{}{
	"proc":       {},
	"sysfs":      {},
	"tmpfs":      {},
	"devtmpfs":   {},
	"devpts":     {},
	"cgroup":     {},
	"cgroup2":    {},
	"securityfs": {},
	"pstore":     {},
	"debugfs":    {},
	"tracefs":    {},
	"configfs":   {},
	"mqueue":     {},
	"hugetlbfs":  {},
	"fusectl":    {},
	"ramfs":      {},
	"autofs":     {},
}

func shouldSkipFSType(fs string) bool {
	_, ok := skippedTypes[fs]
	return ok
}

func shouldSkipMount(mount string) bool {
	return mount == "/proc" || mount == "/sys" || mount == "/dev" || mount == "/run"
}

func Collect() ([]DiskUsage, error) {
	file, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return []DiskUsage{}, err
	}

	mounts := strings.Split(string(file), "\n")
	usagesList := make([]DiskUsage, 0)

	for _, mount := range mounts {
		fields := strings.Fields(mount)
		if len(fields) < 3 {
			continue
		}
		mountPoint := fields[1]
		fsType := fields[2]
		if !strings.HasPrefix(mountPoint, "/") {
			continue
		}
		if shouldSkipFSType(fsType) || shouldSkipMount(mountPoint) {
			continue
		}
		info, err := ReadDiskCounters(mountPoint)
		if err != nil {
			continue
		}
		usage := ConvertFromCountersToUsage(mountPoint, info)
		usagesList = append(usagesList, usage)
	}
	if len(usagesList) == 0 {
		return []DiskUsage{}, errors.New("no info about any mounts")
	}

	return usagesList, nil
}
