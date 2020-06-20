package params

import (
	"fmt"
	"logviewer/utils"
	"strconv"
)

// SignedParam signed parameter
type SignedParam struct {
	CycLogMsgParam
}

// Style get param style
func (p *SignedParam) Style() ParamStyle {
	return Dec
}

// Build build a param from bytes
func (p *SignedParam) Build(data []byte, size int, name string) int {
	p.name = name
	if len(data) < size {
		p.value = "[Length not enought]"
		return size
	}

	switch size {
	case 1:
		p.value = strconv.Itoa(int(data[0]))
	case 2:
		p.value = strconv.Itoa(int(utils.ReadInt16(data)))
	case 4:
		p.value = strconv.Itoa(int(utils.ReadInt32(data)))
	case 8:
		p.value = strconv.Itoa(int(utils.ReadInt64(data)))
	default:
		p.value = fmt.Sprintf("[Signed %d bytes param not implemented]", size)
	}

	return size
}
