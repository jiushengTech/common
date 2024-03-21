package pathutil

import (
	"strings"
	"time"
)

func GetCurrentDateOnlyPath(path string) string {
	format := time.Now().Format(time.DateOnly)
	if path == "" {
		return format
	}
	return path + "/" + format
}

// GetMinioAccessPath path的第一个/前面是bucket
// 将path的第一个/后面添加当前日期
// 如果path中没有/，则直接在path后面添加当前日期

func GetMinioPath(path string) string {
	trim := strings.Trim(path, "/")
	index := strings.Index(trim, "/")
	if index == -1 {
		path = trim + "/" + GetCurrentDateOnlyPath("")
	} else {
		path = trim + "/" + GetCurrentDateOnlyPath(trim[index+1:])
	}
	return path
}
