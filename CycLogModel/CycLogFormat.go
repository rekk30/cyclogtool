package cyclogmodel

import (
	"encoding/csv"
	"fmt"
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
	style ParamStyle
}

// ReadFormatData read format data from *.csv
func ReadFormatData(filename string) (dicFormatMsg map[uint]FormatMsg) {
	dicFormatMsg = make(map[uint]FormatMsg)
	for index := uint(4294967280); index < uint(4294967285); index++ {
		dicFormatMsg[index] = FormatMsg{}
	}
	dicFormatMsg[0] = FormatMsg{comment: "AccOn"}

	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	fmt.Println("Size = ", len(lines))
	check(err)

	for i := 0; i < len(lines); i++ {
		logID, err := strconv.ParseUint(lines[i][0], 10, 64)
		check(err)

		paramCount, err := strconv.ParseInt(lines[i][2], 10, 64)
		check(err)

		var msg = FormatMsg{
			comment: lines[i][1],
		}

		for p := 1; p <= int(paramCount); p++ {
			size, err := strconv.ParseInt(lines[i+p][0], 10, 64)
			check(err)

			style, err := strconv.ParseInt(lines[i+p][2], 10, 64)

			var childMsg = FormatMsgParam{
				size:  int(size),
				name:  lines[i+p][1],
				style: ParamStyle(style),
			}
			msg.child = append(msg.child, childMsg)
		}

		dicFormatMsg[uint(logID)] = msg
		i += int(paramCount)
	}

	return
}
