package zap

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jiushengTech/common/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	log  *zap.Logger
	Sync func() error
}

func NewZapLogger(c *conf.ZapConf) *Logger {
	logger := Logger{}
	cores := logger.GetZapCores(c)
	options := []zap.Option{
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCallerSkip(2),
	}
	if c.Model == "dev" {
		options = append(options, zap.Development())
	}
	// 创建 Logger
	zapLogger := zap.New(zapcore.NewTee(cores...), options...)
	return &Logger{log: zapLogger, Sync: zapLogger.Sync}
}

// GetEncoder 获取 zap core.Encoder
// Author Samsaralc
func (z *Logger) GetEncoder(c *conf.ZapConf) zapcore.Encoder {
	if c.Format == "json" {
		return zapcore.NewJSONEncoder(z.GetEncoderConfig(c))
	}
	return zapcore.NewConsoleEncoder(z.GetEncoderConfig(c))
}

// GetEncoderConfig 获取zap core.EncoderConfig
// Author Samsaralc
func (z *Logger) GetEncoderConfig(c *conf.ZapConf) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  c.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    GetZapEncodeLevel(c.EncodeLevel),
		EncodeTime:     z.CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// GetEncoderCore 获取Encoder的 zap core.Core
// Author Samsaralc
func (z *Logger) GetEncoderCore(c *conf.ZapConf, l zapcore.Level, level zap.LevelEnablerFunc) zapcore.Core {
	writer := z.GetWriteSyncer(c, l.String()) // 日志分割
	return zapcore.NewCore(z.GetEncoder(c), writer, level)
}

// CustomTimeEncoder 自定义日志输出时间格式
// Author Samsaralc
func (z *Logger) CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("2006/01/02 - 15:04:05.000"))
}

// GetZapCores 根据配置文件的Level获取 []zap core.Core
// Author Samsaralc
func (z *Logger) GetZapCores(c *conf.ZapConf) []zapcore.Core {
	cores := make([]zapcore.Core, 0, 7)
	for level := TransportLevel(c.Level); level <= zapcore.FatalLevel; level++ {
		cores = append(cores, z.GetEncoderCore(c, level, GetLevelPriority(level)))
	}
	return cores
}

// GetWriteSyncer 创建日志写入器并设置最大文件大小
// Author Samsaralc
func (z *Logger) GetWriteSyncer(c *conf.ZapConf, level string) zapcore.WriteSyncer {
	logPath := filepath.Join(c.Director, time.Now().Format("2006-01"))
	err := os.MkdirAll(logPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	logFileName := time.Now().Format("02") + "-" + level + ".log"
	logFilePath := filepath.Join(logPath, logFileName)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFilePath,       // 日志文件路径
		MaxSize:    int(c.MaxSize),    // 日志文件的最大大小（以 MB 为单位）
		MaxBackups: int(c.MaxBackups), // 保留的旧日志文件的最大个数
		MaxAge:     int(c.MaxAge),     // 保留的旧日志文件的最大天数
		Compress:   c.Compress,        // 是否压缩旧的日志文件
	}

	// 是否开启控制台输出
	if c.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberjackLogger))
	}
	return zapcore.AddSync(lumberjackLogger)
}

// Log 实现log接口
// Author Samsaralc
func (z *Logger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		z.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}
	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}
	switch level {
	case log.LevelDebug:
		z.log.Debug("", data...)
	case log.LevelInfo:
		z.log.Info("", data...)
	case log.LevelWarn:
		z.log.Warn("", data...)
	case log.LevelError:
		z.log.Error("", data...)
	case log.LevelFatal:
		z.log.Fatal("", data...)
	}
	return nil
}

// GetLevelPriority 根据 zapcore.Level 获取 zap.LevelEnablerFunc
// Author Samsaralc
func GetLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	switch level {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool { // 日志级别
			return level == zap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool { // 警告级别
			return level == zap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool { // 错误级别
			return level == zap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool { // dpanic级别
			return level == zap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool { // panic级别
			return level == zap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool { // 终止级别
			return level == zap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	}
}

// GetZapEncodeLevel 根据 EncodeLevel 返回 zapcore.LevelEncoder
// Author Samsaralc
func GetZapEncodeLevel(encodeLevel string) zapcore.LevelEncoder {
	switch {
	case encodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case encodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case encodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case encodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

// TransportLevel 根据字符串转化为 zapcore.Level
// Author Samsaralc
func TransportLevel(level string) zapcore.Level {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.WarnLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}
