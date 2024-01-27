package log

import (
	"context"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/jiushengTech/common/conf"
	"github.com/jiushengTech/common/log/zap"
	"github.com/jiushengTech/common/utils/file"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var log *klog.Helper

func init() {
	// 构建文件路径
	path := file.CurrentPath()
	filePath := filepath.Join(path, "config.yaml")
	readFile, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var config conf.Config
	err = yaml.Unmarshal(readFile, &config)
	if err != nil {
		panic(err)
	}
	logger := zap.NewZapLogger(config.ZapConf)
	log = klog.NewHelper(logger)
}

func WithContext(ctx context.Context) *klog.Helper {
	return log.WithContext(ctx)
}

func Debug(a ...interface{}) {
	log.Debug(a...)
}

func Debugf(format string, a ...interface{}) {
	log.Debugf(format, a...)
}

func Debugw(keyvals ...interface{}) {
	log.Debugw(keyvals)
}

func Info(a ...interface{}) {
	log.Info(a...)
}

func Infof(format string, a ...interface{}) {
	log.Infof(format, a...)
}

func Infow(keyvals ...interface{}) {
	log.Infow(keyvals...)
}

func Warn(a ...interface{}) {
	log.Warn(a...)
}

func Warnf(format string, a ...interface{}) {
	log.Warnf(format, a...)
}

func Warnw(keyvals ...interface{}) {
	log.Warnw(keyvals...)
}

func Error(a ...interface{}) {
	log.Error(a...)
}

func Errorf(format string, a ...interface{}) {
	log.Errorf(format, a...)
}

func Errorw(keyvals ...interface{}) {
	log.Errorw(keyvals...)
}

func Fatal(a ...interface{}) {
	log.Fatal(a...)
}

func Fatalf(format string, a ...interface{}) {
	log.Fatalf(format, a...)
}

func Fatalw(keyvals ...interface{}) {
	log.Fatalw(keyvals...)
}
