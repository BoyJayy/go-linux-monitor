package disk

import (
	"errors"
	"os"
	"strings"
)

var availableTypes = [8]string{"ext4", "xfs", "btrfs", "f2fs", "vfat", "exfat", "ntfs", "zfs"}
var availableMounts = [5]string{"/", "/boot", "/home", "/mnt/", "/media/"}

func checkFsType(fs string) bool {
	ok := false
	for _, val := range availableTypes {
		if val == fs {
			ok = true
			break
		}
	}
	return ok
}
func checkMount(mount string) bool {
	ok := false
	for _, val := range availableMounts {
		if val == "/" {
			if mount == "/" {
				ok = true
				break
			} else {
				continue
			}
		}
		if string(val[len(val)-1]) != "/" {
			if mount == val {
				ok = true
			} else {
				continue
			}
		}
		if strings.HasPrefix(mount, val) {
			ok = true
			break
		}
	}
	return ok
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
		if !strings.HasPrefix(fields[1], "/") {
			continue
		}
		flag1, flag2 := checkFsType(fields[2]), checkMount(fields[1])
		if !flag1 || !flag2 {
			continue
		}
		info, err := ReadDiskCounters(fields[1])
		if err != nil {
			//return []DiskUsage{}, err
			continue
		}
		usage := ConvertFromCountersToUsage(fields[1], info)
		usagesList = append(usagesList, usage)
	}
	if len(usagesList) == 0 {
		return []DiskUsage{}, errors.New("No info about any mounts")
	}
	return usagesList, nil
}
