package params

import (
	"fmt"
	"strconv"
)

// FloatParam float parameter
type FloatParam struct {
	CycLogMsgParam
}

// Style get param style
func (p *FloatParam) Style() ParamStyle {
	return Float
}

// Build build a param from bytes
func (p *FloatParam) Build(data []byte, size int, name string) int {
	p.name = name
	if len(data) < size {
		p.value = "[Length not enought]"
		return size
	}

	val, err := strconv.ParseFloat(string(data), size*8)
	if err == nil {
		p.value = fmt.Sprintf("%f", val)
	} else {
		p.value = "[Error parsing float value]"
	}

	return size
}
