package zap

import (
	"context"
	"errors"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/jiushengTech/common/log/zap/conf"
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
		Model:         "dev",                        // 开发模式配置
		Level:         "debug",                      // 日志级别设置为 debug（捕获 debug、info、warn、error 等）
		Format:        "console",                    // 日志输出格式（console 或 JSON）
		Director:      "logs",                       // 日志文件存储目录
		EncodeLevel:   "LowercaseColorLevelEncoder", // 使用彩色小写级别名称在日志中
		StacktraceKey: "stack",                      // 堆栈跟踪信息的 JSON 键名
		MaxAge:        0,                            // 保留旧日志文件的最大天数（0 表示无限制）
		ShowLine:      true,                         // 显示日志打印所在的行号
		LogInConsole:  true,                         // 是否在控制台输出日志
		MaxSize:       10,                           // 每个日志文件的最大大小（单位：MB）
		Compress:      false,                        // 是否压缩/归档旧日志文件
		MaxBackups:    10,                           // 保留的旧日志文件的最大数量
	}
	InitLog(&c) // 使用默认配置初始化日志
}

// InitLog initLog 初始化日志系统。
// 如果c为nil，函数会panic。
func InitLog(c *conf.ZapConf) {
	if c == nil {
		panic(errors.New("ZapConf cannot be nil"))
	}
	logger := NewZapLogger(c)
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
