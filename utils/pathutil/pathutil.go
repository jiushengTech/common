package pathutil

import (
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func GetCurrentPath() string {
	var absPath string
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		absPath = filepath.Dir(filename)
	}

	return absPath
}

func GetCurrentDateOnlyAsDir(path string) string {
	trim := strings.Trim(path, "/")
	format := time.Now().Format(time.DateOnly)
	if trim == "" {
		return format
	}
	return path + "/" + format
}
