package asynq

import (
	"fmt"
	z "github.com/jiushengTech/common/log/zap"
	"github.com/jiushengTech/common/log/zap/conf"
	"go.uber.org/zap"

	"github.com/hibiken/asynq"
)

const (
	logKey = "asynq"
)

type logger struct {
	logger *zap.Logger
}

func newLogger() asynq.Logger {
	l := z.NewZapLogger(&conf.ZapConf{
		Model:         "dev",                        // 开发模式配置
		Level:         "debug",                      // 日志级别设置为 debug（捕获 debug、info、warn、error 等）
		Format:        "console",                    // 日志输出格式（console 或 JSON）
		Director:      "logs",                       // 日志文件存储目录
		EncodeLevel:   "LowercaseColorLevelEncoder", // 使用彩色小写级别名称在日志中
		StacktraceKey: "stack",                      // 堆栈跟踪信息的 JSON 键名
		MaxAge:        0,                            // 保留旧日志文件的最大天数（0 表示无限制）
		AddCaller:     true,                         // 显示日志打印所在的行号
		AddCallerSkip: 1,                            // 跳过调用栈的行数
		LogInConsole:  true,                         // 是否在控制台输出日志
		MaxSize:       10,                           // 每个日志文件的最大大小（单位：MB）
		Compress:      true,                         // 是否压缩/归档旧日志文件
		MaxBackups:    10,                           // 保留的旧日志文件的最大数量
		TimeRotation:  z.RotateHourly,               // 时间轮转类型: "0:minute", "1:hour" 或 "2:day"
	})
	return logger{
		logger: l,
	}
}

func (l logger) Debug(args ...interface{}) {
	l.logger.Sugar().Debugln(logKey, fmt.Sprint(args...))
}

func (l logger) Info(args ...interface{}) {
	l.logger.Sugar().Infoln(logKey, fmt.Sprint(args...))
}

func (l logger) Warn(args ...interface{}) {
	l.logger.Sugar().Warnln(logKey, fmt.Sprint(args...))
}

func (l logger) Error(args ...interface{}) {
	l.logger.Sugar().Errorln(logKey, fmt.Sprint(args...))
}

func (l logger) Fatal(args ...interface{}) {
	l.logger.Sugar().Fatalln(logKey, fmt.Sprint(args...))
}
