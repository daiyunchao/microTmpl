package lib

import "time"

func GetNowOffTimeStamp() int64 {
	return time.Now().UnixMicro()
}
