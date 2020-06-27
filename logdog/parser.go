package logdog

import (
	"logviewer/cyclogmodel"
	"logviewer/loginfo"
	"logviewer/utils"
	"os"
	"regexp"
	"strings"
)

// Parser parser for logdog format
type Parser struct {
	model *cyclogmodel.CycLogModel
}

// MakeLogDogParser make logdog
func MakeLogDogParser(m *cyclogmodel.CycLogModel) loginfo.IParser {
	return Parser{m}
}

// Format Parser format messages
func (b Parser) Format(rawMsg cyclogmodel.CycLogMsg) string {
	str := rawMsg.Str
	re := regexp.MustCompile(`%.?`)
	for _, param := range rawMsg.Child {
		str = strings.Replace(str, re.FindString(str), param.Value(), 1)
	}
	return str
}

// Parse Parser parsing function
func (b Parser) Parse(filename string, format map[uint]cyclogmodel.FormatMsg) (logs []loginfo.LogInfo) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	cyclogs := b.model.CycLogReadFile(file, format)

	for _, log := range cyclogs {
		msg := loginfo.LogInfo{
			Text:   b.Format(log),
			Source: utils.ExtractFilename(filename),
			Date:   log.Date,
			LogID:  log.LogID,
		}
		logs = append(logs, msg)
	}

	return
}
