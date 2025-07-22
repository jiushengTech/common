package logger

import (
	"fmt"
	klog "github.com/go-kratos/kratos/v2/log"
	z "github.com/jiushengTech/common/log/zap"
	"github.com/jiushengTech/common/log/zap/conf"
	"go.uber.org/zap"
	"sync"
)

var _ klog.Logger = (*Logger)(nil)

// Logger 实现了klog.Logger接口
type Logger struct {
	logger *zap.Logger
}

var (
	Log  *klog.Helper
	once sync.Once
)

func Init() {
	once.Do(func() {
		logger := z.NewZapLogger(&conf.ZapConf{
			Model:         "dev",                        // 开发模式配置
			Level:         "debug",                      // 日志级别设置为 debug（捕获 debug、info、warn、error 等）
			Format:        "console",                    // 日志输出格式（console 或 JSON）
			Director:      "logs",                       // 日志文件存储目录
			EncodeLevel:   "LowercaseColorLevelEncoder", // 使用彩色小写级别名称在日志中
			StacktraceKey: "stack",                      // 堆栈跟踪信息的 JSON 键名
			MaxAge:        0,                            // 保留旧日志文件的最大天数（0 表示无限制）
			AddCaller:     true,                         // 显示日志打印所在的行号
			AddCallerSkip: 2,                            // 跳过调用栈的行数
			LogInConsole:  true,                         // 是否在控制台输出日志
			MaxSize:       10,                           // 每个日志文件的最大大小（单位：MB）
			Compress:      true,                         // 是否压缩/归档旧日志文件
			MaxBackups:    10,                           // 保留的旧日志文件的最大数量
			TimeRotation:  z.RotateHourly,               // 时间轮转类型: "0:minute", "1:hour" 或 "2:day"
		})
		l := &Logger{
			logger: logger,
		}
		Log = klog.NewHelper(l)
	})

}

func NewLogger(c *conf.ZapConf) *Logger {
	logger := z.NewZapLogger(c)
	l := &Logger{
		logger: logger,
	}
	return l
}

// Log 实现klog.Logger接口的Log方法
func (z *Logger) Log(level klog.Level, keyvals ...interface{}) error {
	// 验证输入参数
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		z.logger.Warn("Keyvalues must appear in pairs", zap.Any("keyvals", keyvals))
		return nil
	}

	// 构建日志字段
	fields := make([]zap.Field, 0, len(keyvals)/2)
	for i := 0; i < len(keyvals); i += 2 {
		key := fmt.Sprint(keyvals[i])
		fields = append(fields, zap.Any(key, keyvals[i+1]))
	}

	// 根据日志级别记录日志
	switch level {
	case klog.LevelDebug:
		z.logger.Debug("", fields...)
	case klog.LevelInfo:
		z.logger.Info("", fields...)
	case klog.LevelWarn:
		z.logger.Warn("", fields...)
	case klog.LevelError:
		z.logger.Error("", fields...)
	case klog.LevelFatal:
		z.logger.Fatal("", fields...)
	}

	return nil
}
