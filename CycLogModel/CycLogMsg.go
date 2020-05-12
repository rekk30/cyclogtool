package cyclogmodel

import "time"

// CycLogMsg msg struct
type CycLogMsg struct {
	child     []CycLogMsgParam
	tickCount int
	date      time.Time
	logID     uint
	str       string
	accCount  int
}

// ParamStyle msg param type
type ParamStyle int

const (
	eUDec ParamStyle = iota
	eUHex
	eDec
	eStruct
	eEnum
	eNd
	eStr
	eFloat
	eDump
)

// CycLogMsgParam parameter for cyclog msg
type CycLogMsgParam struct {
	str   string
	style ParamStyle
	value string
}
