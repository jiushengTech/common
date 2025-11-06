package pathutil

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// GetCurrentPath 返回调用者所在文件的目录路径
func GetCurrentPath() string {
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		return filepath.Dir(filename)
	}
	return ""
}

// GetProjectRoot 返回项目根目录（向上查找 go.mod）
func GetProjectRoot() string {
	dir := GetCurrentPath()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// GetCurrentDateOnlyAsDir 在指定路径下拼接当前日期目录 (YYYY-MM-DD)
func GetCurrentDateOnlyAsDir(basePath string) string {
	format := time.Now().Format(time.DateOnly)
	if basePath == "" {
		return format
	}
	return filepath.Join(basePath, format)
}

// GetCurrentDateOnlyAsDirSlash 在输入目录后拼接日期目录（格式：YYYY-MM-DD）
// 并始终使用 "/" 作为分隔符（适合对象存储路径、URL 等场景）
func GetCurrentDateOnlyAsDirSlash(basePath string) string {
	date := time.Now().Format(time.DateOnly)
	if basePath == "" {
		return date
	}

	// 去掉多余的 /
	basePath = strings.TrimRight(basePath, "/")
	return basePath + "/" + date
}

// GetDateDirByLayout 使用自定义时间格式拼接目录，例如 2025/11/06
func GetDateDirByLayout(basePath string, layout string) string {
	format := time.Now().Format(layout)
	if basePath == "" {
		return format
	}
	return filepath.Join(basePath, format)
}

// EnsureDir 确保目录存在，不存在则自动创建
func EnsureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// PathExists 判断路径是否存在
func PathExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// JoinURLPath 用于拼接 URL 或对象存储路径（始终使用 /）
func JoinURLPath(base string, parts ...string) string {
	all := append([]string{strings.TrimRight(base, "/")}, parts...)
	return path.Join(all...)
}

// GetTimestampedFilename 生成带时间戳的文件名（不含目录）
func GetTimestampedFilename(prefix, ext string) string {
	t := time.Now().Format("20060102_150405")
	if ext != "" && !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return prefix + "_" + t + ext
}
