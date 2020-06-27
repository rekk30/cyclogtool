package utils

import (
	"os"
	"path/filepath"
	"regexp"
)

// Check Check for error
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// ExtractFilename get filename from string
func ExtractFilename(path string) string {
	re := regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)

	match := re.FindStringSubmatch(path)
	return match[2]
}

// FindFile find a file
func FindFile(src string, filename string) (file string, err error) {
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if matched, err := filepath.Match(filename, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			file = path
		}
		return nil
	})

	return
}
