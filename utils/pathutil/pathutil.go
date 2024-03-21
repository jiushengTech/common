package pathutil

import (
	"time"
)

func GetCurrentDateOnlyPath(path string) string {
	format := time.Now().Format(time.DateOnly)
	if path == "" {
		return format
	}
	return format + "/" + path
}
