package params

// ParamStyle msg param type
type ParamStyle int

const (
	UDec ParamStyle = iota
	UHex
	Dec = iota + 1
	Struct
	Enum
	Nd
	Str
	Float
	Dump
)

// CycLogMsgParamIntf interface for parameter
type CycLogMsgParamIntf interface {
	Build([]byte, int, string) int
	Style() ParamStyle
	Name() string
	Value() string
}

// CycLogMsgParam parameter for cyclog msg
type CycLogMsgParam struct {
	CycLogMsgParamIntf
	name  string
	style ParamStyle
	value string
}

// Name return name
func (p CycLogMsgParam) Name() string {
	return p.name
}

// Value return value
func (p CycLogMsgParam) Value() string {
	return p.value
}
