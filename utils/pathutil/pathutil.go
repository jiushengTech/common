package pathutil

import (
	"path/filepath"
	"runtime"
	"strings"
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
	trim := strings.Trim(path, "/")
	format := time.Now().Format(time.DateOnly)
	if trim == "" {
		return format
	}
	return path + "/" + format
}
