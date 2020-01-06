package time

import (
	"os"
	"time"
)

// GetFileModifiedTime returns the last modified time of a file.
// If there is an error, it will return 0
func GetFileModifiedTime(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return fileInfo.ModTime().UnixNano()
}

// Now returns current time as a int64 value.
func Now() int64 {
	return time.Now().UnixNano()
}