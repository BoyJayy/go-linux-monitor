package system

import (
	"time"
)

func CollectTimestamp(interval time.Duration) time.Time {
	time.Sleep(interval)
	return time.Now()
}
