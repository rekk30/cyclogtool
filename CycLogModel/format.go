package cyclogmodel

import (
	"encoding/csv"
	"fmt"
	"logviewer/cyclogmodel/params"
	"os"
	"strconv"
)

// FormatMsg format data struct
type FormatMsg struct {
	comment string
	child   []FormatMsgParam
}

// FormatMsgParam parameter for msg
type FormatMsgParam struct {
	size  int
	name  string
	style params.ParamStyle
}

// ReadFormatData read format data from *.csv
func ReadFormatData(filename string) (dicFormatMsg map[uint]FormatMsg) {
	dicFormatMsg = make(map[uint]FormatMsg)
	for index := uint(4294967280); index < uint(4294967285); index++ {
		dicFormatMsg[index] = FormatMsg{}
	}
	dicFormatMsg[0] = FormatMsg{comment: "AccOn"}

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	fmt.Println("Size = ", len(lines))
	Check(err)

	for i := 0; i < len(lines); i++ {
		logID, err := strconv.ParseUint(lines[i][0], 10, 64)
		Check(err)

		paramCount, err := strconv.ParseInt(lines[i][2], 10, 64)
		Check(err)

		var msg = FormatMsg{
			comment: lines[i][1],
		}

		for p := 1; p <= int(paramCount); p++ {
			size, err := strconv.ParseInt(lines[i+p][0], 10, 64)
			Check(err)

			style, err := strconv.ParseInt(lines[i+p][2], 10, 64)

			var childMsg = FormatMsgParam{
				size:  int(size),
				name:  lines[i+p][1],
				style: params.ParamStyle(style),
			}
			msg.child = append(msg.child, childMsg)
		}

		dicFormatMsg[uint(logID)] = msg
		i += int(paramCount)
	}

	return
}

// ParseMsg parse msg
func ParseMsg(logID uint, data []byte, format map[uint]FormatMsg) (msg CycLogMsg) {
	msg.LogID = logID
	formatMsg, isFound := format[logID]
	if !isFound {
		msg.Str = fmt.Sprintf("Unknown Data(0x%x)", logID)
		return
	}

	msg.Str = formatMsg.comment
	size := len(data)
	pos := 0

	for _, child := range formatMsg.child {
		if size <= 0 {
			break
		}

		var param params.CycLogMsgParamIntf

		switch child.style {
		case params.UDec:
			param = new(params.UnsignedParam)
		case params.UHex:
			param = new(params.HexParam)
		case params.Dec:
			param = new(params.SignedParam)
		case params.Str:
			param = new(params.StringParam)
		case params.Float:
			param = new(params.FloatParam)
		default:
			fmt.Printf("Param style %d not implemented\n", child.style)
		}

		if param == nil {
			fmt.Println("nil")
		}

		read := param.Build(data[pos:], child.size, child.name)
		fmt.Println("Param name = ", param.Name(), " Param value = ", param.Value())

		pos += read
		size -= read

		msg.Child = append(msg.Child, param)
	}

	return
}

// func setData(logID uint, data []byte, format map[uint]FormatMsg) (str string, params []CycLogMsgParam) {
