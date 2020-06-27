package loginfo

import (
	"logviewer/cyclogmodel"
	"time"
)

// LogInfo parsed log struct
type LogInfo struct {
	Pid      int
	Tid      int
	Source   string
	LogLevel string
	ProcName string
	Text     string
	Date     time.Time
	LogID    uint
}

// IParser logfile parser interface
type IParser interface {
	Format(cyclogmodel.CycLogMsg) string
	Parse(string, map[uint]cyclogmodel.FormatMsg) []LogInfo
}
