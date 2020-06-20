package params

import (
	"fmt"
	"logviewer/utils"
	"strconv"
)

// UnsignedParam unsigned parameter
type UnsignedParam struct {
	CycLogMsgParam
}

// Style get param style
func (p *UnsignedParam) Style() ParamStyle {
	return UDec
}

// Build build a param from bytes
func (p *UnsignedParam) Build(data []byte, size int, name string) int {
	p.CycLogMsgParam.name = name
	if len(data) < size {
		p.value = "[Length not enought]"
		return size
	}

	switch size {
	case 1:
		p.value = strconv.Itoa(int(data[0]))
	case 2:
		p.value = strconv.Itoa(int(utils.ReadUInt16(data)))
	case 4:
		p.value = strconv.Itoa(int(utils.ReadUInt32(data)))
	case 8:
		p.value = strconv.Itoa(int(utils.ReadUInt64(data)))
	default:
		p.value = fmt.Sprintf("[Unsigned %d bytes param not implemented]", size)
	}

	return size
}
