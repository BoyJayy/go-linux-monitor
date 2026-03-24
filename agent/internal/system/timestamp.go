package system

import (
	"time"
)

func Collect(interval time.Duration) time.Time {
	time.Sleep(interval)
	return time.Now()
}
