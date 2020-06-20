package main

import (
	"fmt"
	"io"
	cyclogmodel "logviewer/CycLogModel"
	logdog "logviewer/LogDog"
	"logviewer/loginfo"
	"logviewer/utils"
	"os"
	"sort"
)

// FormatDir logdog format folder
const FormatDir = "logdogformat"

// LogFileExt log file pattern
const LogFileExt = "*.LOG"

func writeLogs(logs []loginfo.LogInfo, r io.Writer) int {
	count := 0
	for _, log := range logs {
		timeString := fmt.Sprintf("%02d:%02d:%02d.%d",
			log.Date.Hour(), log.Date.Minute(),
			log.Date.Second(), log.Date.Nanosecond())
		logIDString := fmt.Sprintf("[0x%04x] ", log.LogID)

		io.WriteString(r, timeString+" ")
		io.WriteString(r, log.Source+" ")
		io.WriteString(r, logIDString)
		io.WriteString(r, log.Text)
		io.WriteString(r, "\n")

		count++
	}

	return count
}

func main() {
	if len(os.Args) < 2 {
		panic("Specify log archive")
	}

	file, err := os.Open(os.Args[1])
	utils.Check(err)
	defer file.Close()
	logFolder := file.Name() + ".ext"

	var outputFile *os.File
	outputFile, err = os.Create("logs.txt")
	utils.Check(err)
	defer outputFile.Close()

	utils.Untar(logFolder, file)

	model := cyclogmodel.MakeCycLogModel()
	parser := logdog.MakeLogDogParser(&model)

	files := utils.GetFiles(logFolder, LogFileExt)

	var logs []loginfo.LogInfo

	for _, file := range files {
		filename := utils.ExtractFilename(file)

		format := cyclogmodel.ReadFormatData(FormatDir + "/" + filename + ".csv")
		fmt.Println(len(format))
		fmt.Println("Processing " + filename)

		logs = append(logs, parser.Parse(file, format)...)
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].Date.Before(logs[j].Date)
	})

	logsWritten := writeLogs(logs, outputFile)
	fmt.Printf("Written %d logs\n", logsWritten)
}
