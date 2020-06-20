package params

// StringParam string parameter
type StringParam struct {
	CycLogMsgParam
}

// Style get param style
func (p *StringParam) Style() ParamStyle {
	return Str
}

// Build build a param from bytes
func (p *StringParam) Build(data []byte, size int, name string) (read int) {
	p.name = name
	if len(data) < size {
		read = 0
		p.value = "[Length not enought]"
		return
	}

	p.value = string(data[:size])
	read = size

	return
}
