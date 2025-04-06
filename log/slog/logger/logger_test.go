package logger

import (
	"fmt"
	"testing"
)

func TestLogger(t *testing.T) {
	opts := DefaultOptions()
	opts.Format = "json"
	opts.FilePrefix = "myapp"
	opts.EnableFile = true
	opts.EnableStdout = true

	logger, err := NewSlogLogger(opts)
	if err != nil {
		fmt.Println("初始化日志失败:", err)
	}

	logger.Info("应用启动", "env", "production")
	logger.Debug("调试日志", "user", "admin")

}
