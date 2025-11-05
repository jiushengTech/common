package pathutil

import (
	"path/filepath"
	"runtime"
	"time"
)

func GetCurrentPath() string {
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		return filepath.Dir(filename)
	}
	return ""
}

func GetCurrentDateOnlyAsDir(path string) string {
	format := time.Now().Format(time.DateOnly)
	if path == "" {
		return format
	}
	return filepath.Join(path, format)
}
