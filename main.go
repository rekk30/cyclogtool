package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"logviewer/config"
	"logviewer/cyclogmodel"
	"logviewer/logdog"
	"logviewer/loginfo"
	"logviewer/utils"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	// LogFileExt log file pattern
	LogFileExt = "*.LOG"
	// Version tool version
	Version = "0.0.1"
	// Purple color
	Purple = "\033[1;34m"
	// Normal color
	Normal = "\033[0m"
)

func writeLogs(logs []loginfo.LogInfo, r io.Writer, conf config.Config) int {
	count := 0
	for _, log := range logs {
		str := conf.Format
		logIDString := fmt.Sprintf("[0x%04x] ", log.LogID)

		str = strings.ReplaceAll(str, "%y", fmt.Sprintf("%d", log.Date.Year()))
		str = strings.ReplaceAll(str, "%k", fmt.Sprintf("%d", log.Date.Month()))
		str = strings.ReplaceAll(str, "%d", fmt.Sprintf("%d", log.Date.Day()))
		str = strings.ReplaceAll(str, "%h", fmt.Sprintf("%02d", log.Date.Hour()))
		str = strings.ReplaceAll(str, "%m", fmt.Sprintf("%02d", log.Date.Minute()))
		str = strings.ReplaceAll(str, "%s", fmt.Sprintf("%02d", log.Date.Second()))
		str = strings.ReplaceAll(str, "%n", fmt.Sprintf("%09d", log.Date.Nanosecond()))
		str = strings.ReplaceAll(str, "%i", logIDString)
		str = strings.ReplaceAll(str, "%t", log.Text)
		str = strings.ReplaceAll(str, "%x", log.Source)
		str = strings.ReplaceAll(str, "%p", string(log.Pid))
		str = strings.ReplaceAll(str, "%z", string(log.Tid))
		str = strings.ReplaceAll(str, "%l", log.LogLevel)

		io.WriteString(r, str+"\n")

		count++
	}

	return count
}

func encrypted(path string, conf config.Config) (decryptedData []byte, err error) {
	name := strings.TrimSuffix(
		strings.TrimSuffix(path, filepath.Ext(path)),
		"_enckey")

	archive := name + ".enc"
	key := name + "_enckey.enc"

	factoryKeyData, err := ioutil.ReadFile(conf.Key.Factory)
	if err != nil {
		return
	}
	privateKeyData, err := ioutil.ReadFile(conf.Key.Installed)
	if err != nil {
		return
	}
	factoryKey, err := utils.BytesToPrivateKey(factoryKeyData)
	if err != nil {
		return
	}
	privateKey, err := utils.BytesToPrivateKey(privateKeyData)
	if err != nil {
		return
	}

	keyData, err := ioutil.ReadFile(key)
	if err != nil {
		return
	}
	decryptedKey, err := utils.RsaDecrypt(keyData, privateKey, factoryKey)
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(archive)
	if err != nil {
		return
	}

	hx := hex.EncodeToString(decryptedKey)
	decryptedData, err = utils.Aes256Decrypt(data, hx[0:64], hx[64:96])
	return
}

func parseArg(args string, conf config.Config) (logFolder string, err error) {
	path, err := os.Stat(args)

	if err != nil {
		return
	}

	isArchive := false
	var archiveReader io.Reader

	switch mode := path.Mode(); {
	case mode.IsDir():
		logFolder = path.Name()
	case mode.IsRegular():
		switch filepath.Ext(path.Name()) {
		case ".enc":
			var data []byte
			if data, err = encrypted(path.Name(), conf); err != nil {
				return
			}

			isArchive = true
			archiveReader = bytes.NewReader(data)

		case ".gz":
			var file *os.File
			file, err = os.Open(path.Name())
			if err != nil {
				return
			}

			isArchive = true
			archiveReader = file
		}
	}

	if isArchive {
		if !conf.UseTemp {
			logFolder = path.Name() + ".ext"
		} else {
			logFolder, err = ioutil.TempDir("", path.Name())
			if err != nil {
				return
			}
		}

		utils.Untar(logFolder, archiveReader)
	}

	return
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%sGoLogdog v%s%s\n", Purple, Version, Normal)
		fmt.Fprintf(os.Stderr, "Usage: [ARGS] (tar.gz|enc|/)\n")
		fmt.Fprintf(os.Stderr, "Format:\n")
		fmt.Fprintf(os.Stderr, "%%y   year\n")
		fmt.Fprintf(os.Stderr, "%%k   month\n")
		fmt.Fprintf(os.Stderr, "%%d   day\n")
		fmt.Fprintf(os.Stderr, "%%h   hours\n")
		fmt.Fprintf(os.Stderr, "%%m   minutes\n")
		fmt.Fprintf(os.Stderr, "%%s   seconds\n")
		fmt.Fprintf(os.Stderr, "%%n   nanoseconds\n")
		fmt.Fprintf(os.Stderr, "%%i   logID\n")
		fmt.Fprintf(os.Stderr, "%%t   text\n")
		fmt.Fprintf(os.Stderr, "%%x   source\n")
		fmt.Fprintf(os.Stderr, "%%p   pid\n")
		fmt.Fprintf(os.Stderr, "%%z   tid\n")
		fmt.Fprintf(os.Stderr, "%%l   logLevel\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		flag.PrintDefaults()
	}

	configFile := flag.String("c", config.ConfigFile, "Configuration file")
	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	var configuration config.Config
	configuration, err := config.GetConfig(*configFile)

	if err != nil {
		fmt.Println("Can't read configuration, using default")
		configuration = config.GetDefaultConfig()
	}

	logFolder, err := parseArg(os.Args[1], configuration)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(logFolder)

	outputFile, err := os.Create(os.Args[1] + ".log")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	model := cyclogmodel.MakeCycLogModel()
	parser := logdog.MakeLogDogParser(&model)

	var logs []loginfo.LogInfo

	for _, file := range configuration.Target.Logdog {
		// File
		logPath, err := utils.FindFile(logFolder, file+LogFileExt)
		if err != nil {
			fmt.Printf("File \"%s\" not found\n", logFolder+"/"+file+LogFileExt)
			continue
		}

		// csv
		csvFile := configuration.Target.LogdogFormat + "/" + file + ".csv"
		format, err := cyclogmodel.ReadFormatData(csvFile)
		if err != nil {
			fmt.Printf("File \"%s\" not found\n", csvFile)
			continue
		}

		logs = append(logs, parser.Parse(logPath, format)...)
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].Date.Before(logs[j].Date)
	})

	logsWritten := writeLogs(logs, outputFile, configuration)
	fmt.Printf("Written %d logs -> \"%s\"\n", logsWritten, outputFile.Name())
}
