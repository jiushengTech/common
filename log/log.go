package log

import (
	"context"
	"errors"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/jiushengTech/common/conf"
	"github.com/jiushengTech/common/log/zap"
)

var log *klog.Helper

func init() {
	//// 构建文件路径
	//path := fileutil.CurrentPath()
	//filePath := filepath.Join(path, "config.yaml")
	//readFile, err := os.ReadFile(filePath)
	//if err != nil {
	//	panic(err)
	//}
	//var config conf.Config
	//err = yaml.Unmarshal(readFile, &config)
	//if err != nil {
	//	panic(err)
	//}
	// 默认配置
	c := conf.ZapConf{
		Model:         "dev",
		Level:         "debug",
		Format:        "console",
		Director:      "logs",
		EncodeLevel:   "LowercaseColorLevelEncoder",
		StacktraceKey: "stack",
		MaxAge:        0,
		ShowLine:      true,
		LogInConsole:  true,
		MaxSize:       10,
		Compress:      false,
		MaxBackups:    10,
	}
	InitLog(&c) // 使用默认配置初始化日志
}

// initLog 初始化日志系统。
// 如果c为nil，函数会panic。
func InitLog(c *conf.ZapConf) {
	if c == nil {
		panic(errors.New("ZapConf cannot be nil"))
	}
	logger := zap.NewZapLogger(c)
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
