package cyclogmodel

import (
	"logviewer/utils"
	"time"
)

// CycLogTime time struct
type CycLogTime struct {
	tickCount int
	date      time.Time
	accCount  int
	logID     uint
}

func setTime(data []byte) time.Time {
	if len(data) < 16 {
		return time.Time{}
	}
	year := int(utils.ReadInt16(data))
	month := time.Month(utils.ReadInt16(data[2:]))
	day := int(utils.ReadInt16(data[6:]))
	hour := int(utils.ReadInt16(data[8:]))
	minute := int(utils.ReadInt16(data[10:]))
	second := int(utils.ReadInt16(data[12:]))
	millisecond := int(utils.ReadInt16(data[14:]))

	return time.Date(year, month, day, hour, minute, second, millisecond, time.UTC)
}
