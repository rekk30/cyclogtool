package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"logviewer/cyclogmodel"
	"os"
	"path/filepath"
	"strings"
)

// FormatDir logdog format folder
const FormatDir = "logdogformat"

// LogFileExt log file pattern
const LogFileExt = "*.LOG"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func untar(dst string, r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

func getFiles(src string, pattern string) []string {
	var files []string

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			files = append(files, path)
		}
		return nil
	})
	check(err)

	return files
}

// LogInfo parsed log struct
type LogInfo struct {
	source   int
	fullText string
	pid      int
	tid      int
	logLevel string
	procName string
	text     string
}

// IParser logfile parser interface
type IParser interface {
	Parse(string) LogInfo
	ParseLogTime(string, map[uint]cyclogmodel.FormatMsg) LogInfo
}

// LogDogParser parser for logdog format
type LogDogParser struct {
	model *cyclogmodel.CycLogModel
}

func makeLogDogParser(m *cyclogmodel.CycLogModel) IParser {
	return LogDogParser{m}
}

// Parse LogDogParser parsing function
func (b LogDogParser) Parse(filename string) LogInfo {
	_, err := ioutil.ReadFile(filename)
	check(err)

	// b.model.CycLogReadFile(filename)

	return LogInfo{}
}

// ParseLogTime parse logtime.dat file
func (b LogDogParser) ParseLogTime(filename string, format map[uint]cyclogmodel.FormatMsg) LogInfo {
	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	b.model.CycLogReadLogTime(file, format)

	return LogInfo{}
}

func main() {
	if len(os.Args) < 2 {
		panic("Specify log archive")
	}
	file, err := os.Open(os.Args[1])
	check(err)
	logFolder := file.Name() + ".ext"

	untar(logFolder, file)

	// timeLogFile := getFiles(logFolder, "LOGTIME.DAT")
	// if len(timeLogFile) < 1 {
	// 	panic("LOGTIME.DAT not find")
	// }

	model := cyclogmodel.MakeCycLogModel()
	parser := makeLogDogParser(&model)

	// timeLog := parser.ParseLogTime(timeLogFile[0])

	// fmt.Println(timeLog)

	files := getFiles(logFolder, LogFileExt)

	for _, file := range files {
		parts := strings.Split(file, string(os.PathSeparator))
		filename := parts[len(parts)-1]

		parts2 := strings.Split(filename, ".")

		format := cyclogmodel.ReadFormatData(FormatDir + "/" + parts2[0] + ".csv")
		fmt.Println(len(format))
		fmt.Println("Processing " + filename)
		parser.ParseLogTime(file, format)
	}
	// fmt.Scanln()
}
