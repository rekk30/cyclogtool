package cyclogmodel

import (
	"logviewer/utils"
	"os"
	"time"
)

// Check check for error
func Check(e error) {
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

	Date []CycLogTime
}

// MakeCycLogModel model
func MakeCycLogModel() CycLogModel {
	var model CycLogModel

	return model
}

// CycLogReadHeader read header
func (model CycLogModel) CycLogReadHeader(file *os.File) (filePosition int, logType CycLogHeadType) {
	fileStat, err := file.Stat()
	Check(err)

	if fileStat.Size() < 32 {
		panic("File size < 32")
	}

	bytes := make([]byte, 32)
	_, err = file.Read(bytes)
	Check(err)

	filePosition = 0
	logType = CycLogHeadType(string(bytes[:12]))
	if _, ok := cycLogTypes[logType]; !ok {
		panic("Wrong log type")
	}

	model.headWriteAddr[0] = utils.ReadInt32(bytes[12:])
	model.headWriteAddr[1] = utils.ReadInt32(bytes[16:])
	model.headTailAddr = utils.ReadInt32(bytes[20:])
	model.headUpdIdx = bytes[24]
	model.headRecordSize = utils.ReadInt16(bytes[28:])

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
		recordSize = int(utils.ReadInt16(numArray[6:])) - cycLogTypes[logType].size
	} else if logType == eCycLogRecordEx {
		recordSize = int(utils.ReadInt16(numArray[8:])) - cycLogTypes[logType].size
	}

	if recordSize < 0 {
		panic("Record size < 0")
	}

	return
}

func readLogID(logType CycLogHeadType, numArray []byte) (logID uint) {
	logID = 0

	if cycLogTypes[logType].logIDSize == 16 {
		logID = uint(utils.ReadUInt16(numArray[4:]))
		if logType == eCycLogRecord && ((logID & 65520) == 65520) {
			logID = uint(-65536 | int(logID)&int(65535))
		}
	} else {
		logID = uint(utils.ReadUInt32(numArray[4:]))
	}

	return
}

// CycLogReadFile read cyc log file
func (model CycLogModel) CycLogReadFile(file *os.File, format map[uint]FormatMsg) (logMsg []CycLogMsg) {
	fileStat, err := file.Stat()
	Check(err)

	filePosition, logType := model.CycLogReadHeader(file)
	logTypeInfo := cycLogTypes[logType]

	flag := false

	for ok, logID := true, uint(0); ok; ok = (logID != 0) {
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
		Check(err)
		_, err = file.Read(numArray)
		Check(err)

		ticks := int(utils.ReadInt32(numArray))

		logID = readLogID(logType, numArray)

		recordSize := readRecordSize(logType, numArray)

		if !flag || filePosition >= int(model.headWriteAddr[model.headUpdIdx])+recordSize {
			filePosition -= recordSize
			if filePosition < 32 || fileStat.Size() < int64(filePosition) {
				panic("Bad size")
			}

			numArray3 := make([]byte, recordSize)
			_, err := file.Seek(int64(filePosition), 0)
			Check(err)
			_, err = file.Read(numArray3)
			Check(err)

			logMsg1 := ParseMsg(logID, numArray3, format)

			logMsg1.Date = logMsg1.Date.Add(time.Millisecond * time.Duration(ticks))

			logMsg = append(logMsg, logMsg1)
		}
	}

	return
}
