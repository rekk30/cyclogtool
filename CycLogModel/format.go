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
func ReadFormatData(filename string) (dicFormatMsg map[uint]FormatMsg, err error) {
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

	for i := 0; i < len(lines); i++ {
		logID, _ := strconv.ParseUint(lines[i][0], 10, 64)
		// if err == nil {
		// 	continue
		// }

		paramCount, _ := strconv.ParseInt(lines[i][2], 10, 64)
		// if err == nil {
		// 	continue
		// }

		var msg = FormatMsg{
			comment: lines[i][1],
		}

		for p := 1; p <= int(paramCount); p++ {
			size, _ := strconv.ParseInt(lines[i+p][0], 10, 64)
			// if err == nil {
			// 	continue
			// }

			style, _ := strconv.ParseInt(lines[i+p][2], 10, 64)

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
		}

		read := 0
		if param != nil {
			read = param.Build(data[pos:], child.size, child.name)
			msg.Child = append(msg.Child, param)
		}

		pos += read
		size -= read

	}

	return
}

// func setData(logID uint, data []byte, format map[uint]FormatMsg) (str string, params []CycLogMsgParam) {
