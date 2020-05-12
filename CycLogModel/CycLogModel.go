package cyclogmodel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// CycLogHeadType log msg type
type CycLogHeadType string

const (
	eCycLogRecord      CycLogHeadType = "CYCLICLOG201"
	eCycLogRecordEx                   = "CYCLICLOG301"
	eFixCycLogRecord                  = "FIXEDLOG 201"
	eFixCycLogRecordEx                = "FIXEDLOG 301"
)

// CycLogHeadInfo info about header
type CycLogHeadInfo struct {
	size      int
	logIDSize int
}

var cycLogTypes = map[CycLogHeadType]CycLogHeadInfo{
	eCycLogRecord:      {8, 16},
	eCycLogRecordEx:    {12, 32},
	eFixCycLogRecord:   {8, 16},
	eFixCycLogRecordEx: {8, 32},
}

// CycLogModel main cyclog model
type CycLogModel struct {
	headWriteAddr  [2]int32
	headTailAddr   int32
	headUpdIdx     byte
	headRecordSize int16
}

// MakeCycLogModel model
func MakeCycLogModel() CycLogModel {
	var model CycLogModel

	return model
}

func readSome(data []byte, ret interface{}) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func readInt16(data []byte) (ret int16) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func readInt32(data []byte) (ret int32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func readInt64(data []byte) (ret int64) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func readUInt16(data []byte) (ret uint16) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func readUInt32(data []byte) (ret uint32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func readUInt64(data []byte) (ret uint64) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

// CycLogReadHeader read header
func (model CycLogModel) CycLogReadHeader(file *os.File) (filePosition int, logType CycLogHeadType) {
	fileStat, err := file.Stat()
	check(err)

	if fileStat.Size() < 32 {
		panic("File size < 32")
	}

	bytes := make([]byte, 32)
	_, err = file.Read(bytes)
	check(err)

	filePosition = 0
	logType = CycLogHeadType(string(bytes[:12]))
	if _, ok := cycLogTypes[logType]; !ok {
		panic("Wrong log type")
	}

	model.headWriteAddr[0] = readInt32(bytes[12:])
	model.headWriteAddr[1] = readInt32(bytes[16:])
	model.headTailAddr = readInt32(bytes[20:])
	model.headUpdIdx = bytes[24]
	model.headRecordSize = readInt16(bytes[28:])

	// fmt.Println("Head tail addr ", model.headTailAddr)
	// fmt.Println("Head upd idx", model.headUpdIdx)
	// fmt.Println("Head record size", model.headRecordSize)

	// if logType != eFixCycLogRecord && logType != eFixCycLogRecordEx {
	// 	recordSize = int(model.headRecordSize) - model.cycLogTypes[logType].size
	// }

	if model.headUpdIdx < 0 || model.headUpdIdx > 1 {
		panic("Wrong headUpdIdx value")
	}

	filePosition = int(model.headWriteAddr[model.headUpdIdx])

	if filePosition < 32 || int64(filePosition) > fileStat.Size() {
		panic("Wrong file position")
	}
	return
}

func readRecordSize(logType CycLogHeadType, numArray []byte) (recordSize int) {
	recordSize = 0

	if logType == eCycLogRecord {
		recordSize = int(readInt16(numArray[6:])) - cycLogTypes[logType].size
	} else if logType == eCycLogRecordEx {
		recordSize = int(readInt16(numArray[8:])) - cycLogTypes[logType].size
	}

	if recordSize < 0 {
		panic("Record size < 0")
	}

	return
}

func readLogID(logType CycLogHeadType, numArray []byte) (logID uint) {
	logID = 0

	if cycLogTypes[logType].logIDSize == 16 {
		logID = uint(readUInt16(numArray[4:]))
		if logType == eCycLogRecord && ((logID & 65520) == 65520) {
			logID = uint(-65536 | int(logID)&int(65535))
		}
	} else {
		logID = uint(readUInt32(numArray[4:]))
	}

	return
}

func binToUStr(data []byte, pos int, size int) (str string, result bool) {
	result = false
	str = ""
	if len(data) < pos+size {
		return
	}

	switch size {
	case 1:
		str = string(data[pos])
	case 2:
		str = string(readUInt16(data[pos:]))
	case 4:
		str = string(readUInt32(data[pos:]))
	case 8:
		str = string(readUInt64(data[pos:]))
	default:
		panic("Not implemented")
	}

	result = true
	return
}

func setTime(tickCount int) (date time.Time) {

	return
}

func setData(logID uint, data []byte, format map[uint]FormatMsg) (str string, params []CycLogMsgParam) {
	formatMsg, isFound := format[logID]
	if !isFound {
		str = fmt.Sprintf("Unknown Data(0x%x)", logID)
		return
	}

	str = formatMsg.comment
	size := len(data)
	pos := 0
	// fmt.Println("Size = ", size, " childs = ", len(formatMsg.child))

	for _, child := range formatMsg.child {
		if size <= 0 {
			break
		}

		param := CycLogMsgParam{
			str:   child.name,
			style: child.style,
		}

		result := false
		switch child.style {
		case eUDec:
			fallthrough
		case eUHex:
			param.value, result = binToUStr(data, pos, child.size)
			if result {
				pos += child.size
				size -= child.size
			} //TODO error handling
		default:
			// fmt.Println("Not implemented")

		}

		params = append(params, param)
	}

	return
}

// // CycLogReadFile read cyc log file
// func (model CycLogModel) CycLogReadFile(filename string) {

// CycLogReadLogTime read cyc log file
func (model CycLogModel) CycLogReadLogTime(file *os.File, format map[uint]FormatMsg) (logMsg []CycLogMsg) {
	fileStat, err := file.Stat()
	check(err)

	filePosition, logType := model.CycLogReadHeader(file)

	// fmt.Println("File position = ", filePosition)
	// fmt.Println("Log type = ", logType)

	flag := false
	// var timeLogMsg1 CycLogMsg
	// index := -1
	// num1 := 0
	num2 := 0

	// for {
	for ok, logID := true, uint(0); ok; ok = (logID != 0) {
		logTypeInfo := cycLogTypes[logType]
		if filePosition == 32 {
			if model.headTailAddr != 0 {
				flag = true
				filePosition = int(model.headTailAddr)
			} else {
				break
			}
		} else if filePosition < 32 {
			return
		}

		filePosition -= logTypeInfo.size
		if filePosition < 32 {
			panic("File position < 32")
		}

		numArray := make([]byte, logTypeInfo.size)
		_, err := file.Seek(int64(filePosition), 0)
		check(err)
		_, err = file.Read(numArray)
		check(err)

		ticks := readInt32(numArray)
		// fmt.Println("Tick count = ", ticks)

		logID = readLogID(logType, numArray)
		// fmt.Println("logID = ", logID)

		recordSize := readRecordSize(logType, numArray)

		if !flag || filePosition >= int(model.headWriteAddr[model.headUpdIdx])+recordSize {
			filePosition -= recordSize
			if filePosition < 32 || fileStat.Size() < int64(filePosition) {
				panic("Bad size")
			}

			// fmt.Println("LogID = ", logID, " Record size = ", recordSize)
			numArray3 := make([]byte, recordSize)

			_, err := file.Seek(int64(filePosition), 0)
			check(err)
			_, err = file.Read(numArray3)
			check(err)

			logMsg1 := CycLogMsg{
				tickCount: int(ticks),
				logID:     logID,
				accCount:  num2,
			}

			logMsg1.date = setTime(logMsg1.tickCount)
			logMsg1.str, logMsg1.child = setData(logID, numArray3, format)
			logMsg = append(logMsg, logMsg1)
		}
	}
	// }

	return
}
