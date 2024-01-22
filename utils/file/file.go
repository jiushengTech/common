package file

import (
	"path/filepath"
	"runtime"
)

func CurrentPath() string {
	var absPath string
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		absPath = filepath.Dir(filename)
	}

	return absPath
}
