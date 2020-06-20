package cyclogmodel

import (
	"logviewer/cyclogmodel/params"
	"time"
)

// CycLogMsg msg struct
type CycLogMsg struct {
	Child []params.CycLogMsgParamIntf
	Date  time.Time
	LogID uint
	Str   string
	//TODO add error
}
