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

// GetFiles get all files
func GetFiles(src string, pattern string) []string {
	var files []string

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			files = append(files, path)
		}
		return nil
	})
	Check(err)

	return files
}
