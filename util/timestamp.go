package util

import "time"

func TimeStamp() int64 {
	return time.Now().UnixNano()
}
