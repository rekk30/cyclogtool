package params

import (
	"fmt"
	"logviewer/utils"
)

// HexParam hex parameter
type HexParam struct {
	CycLogMsgParam
}

// Style get param style
func (p *HexParam) Style() ParamStyle {
	return UHex
}

// Build build a param from bytes
func (p *HexParam) Build(data []byte, size int, name string) int {
	p.name = name
	if len(data) < size {
		p.value = "[Length not enought]"
		return size
	}

	switch size {
	case 1:
		p.value = fmt.Sprintf("0x%X", int(data[0]))
	case 2:
		p.value = fmt.Sprintf("0x%X", int(utils.ReadUInt16(data)))
	case 4:
		p.value = fmt.Sprintf("0x%X", int(utils.ReadUInt32(data)))
	case 8:
		p.value = fmt.Sprintf("0x%X", int(utils.ReadUInt64(data)))
	default:
		p.value = fmt.Sprintf("[Hex %d bytes param not implemented]", size)
	}

	return size
}
