package zap

import (
	"fmt"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/jiushengTech/common/log/zap/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var _ klog.Logger = (*Logger)(nil)

// Logger 实现了klog.Logger接口
type Logger struct {
	log  *zap.Logger
	Sync func() error
}

// NewZapLogger 创建一个新的zap日志记录器
func NewZapLogger(c *conf.ZapConf) *Logger {
	logger := &Logger{}
	cores := logger.GetZapCores(c)
	options := []zap.Option{
		zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCallerSkip(2),
	}

	if c.Model == "dev" {
		options = append(options, zap.Development())
	}

	zapLogger := zap.New(zapcore.NewTee(cores...), options...)
	return &Logger{log: zapLogger, Sync: zapLogger.Sync}
}

// GetEncoder 获取编码器
func (z *Logger) GetEncoder(c *conf.ZapConf, forConsole bool) zapcore.Encoder {
	encoderConfig := z.GetEncoderConfig(c, forConsole)

	if c.Format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// GetEncoderConfig 获取编码器配置
func (z *Logger) GetEncoderConfig(c *conf.ZapConf, forConsole bool) zapcore.EncoderConfig {
	encConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  c.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeTime:     z.CustomTimeEncoder,
	}

	if forConsole {
		// 控制台输出使用彩色编码器
		encConfig.EncodeLevel = GetZapEncodeLevel(c.EncodeLevel)
	} else {
		// 文件输出使用非彩色编码器
		encConfig.EncodeLevel = getNonColorEncoder(c.EncodeLevel)
	}

	return encConfig
}

// getNonColorEncoder 获取非彩色级别编码器
func getNonColorEncoder(encodeLevel string) zapcore.LevelEncoder {
	switch encodeLevel {
	case "LowercaseColorLevelEncoder", "LowercaseLevelEncoder":
		return zapcore.LowercaseLevelEncoder
	case "CapitalColorLevelEncoder", "CapitalLevelEncoder":
		return zapcore.CapitalLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

// CustomTimeEncoder 自定义日志时间格式
func (z *Logger) CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format(time.DateTime))
}

// GetZapCores 获取所有zap核心组件
func (z *Logger) GetZapCores(c *conf.ZapConf) []zapcore.Core {
	cores := make([]zapcore.Core, 0, 7)
	minLevel := TransportLevel(c.Level)

	// 为每个级别创建Core
	for level := minLevel; level <= zapcore.FatalLevel; level++ {
		levelFunc := GetLevelPriority(level)

		// 添加文件输出Core
		fileCore := z.createFileCore(c, level, levelFunc)
		cores = append(cores, fileCore)

		// 如果需要控制台输出，添加控制台Core
		if c.LogInConsole {
			consoleCore := z.createConsoleCore(c, levelFunc)
			cores = append(cores, consoleCore)
		}
	}

	return cores
}

// createFileCore 创建文件输出的Core
func (z *Logger) createFileCore(c *conf.ZapConf, level zapcore.Level, levelFunc zap.LevelEnablerFunc) zapcore.Core {
	return zapcore.NewCore(
		z.GetEncoder(c, false),
		z.GetWriteSyncer(c, level.String()),
		levelFunc,
	)
}

// createConsoleCore 创建控制台输出的Core
func (z *Logger) createConsoleCore(c *conf.ZapConf, levelFunc zap.LevelEnablerFunc) zapcore.Core {
	return zapcore.NewCore(
		z.GetEncoder(c, true),
		zapcore.AddSync(os.Stdout),
		levelFunc,
	)
}

// GetWriteSyncer 创建文件日志写入器
func (z *Logger) GetWriteSyncer(c *conf.ZapConf, level string) zapcore.WriteSyncer {
	// 创建日志目录
	logPath := filepath.Join(c.Director, time.Now().Format("2006-01"))
	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		panic(fmt.Sprintf("创建日志目录失败: %v", err))
	}

	logFileName := fmt.Sprintf("%02d-%s.log", time.Now().Day(), level)
	logFilePath := filepath.Join(logPath, logFileName)

	// 配置日志文件分割
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    int(c.MaxSize),
		MaxAge:     int(c.MaxAge),
		MaxBackups: int(c.MaxBackups),
		LocalTime:  true,
		Compress:   c.Compress,
	}

	return zapcore.AddSync(lumberjackLogger)
}

// Log 实现klog.Logger接口的Log方法
func (z *Logger) Log(level klog.Level, keyvals ...interface{}) error {
	// 验证输入参数
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		z.log.Warn("Keyvalues must appear in pairs", zap.Any("keyvals", keyvals))
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
		z.log.Debug("", fields...)
	case klog.LevelInfo:
		z.log.Info("", fields...)
	case klog.LevelWarn:
		z.log.Warn("", fields...)
	case klog.LevelError:
		z.log.Error("", fields...)
	case klog.LevelFatal:
		z.log.Fatal("", fields...)
	}

	return nil
}

// GetLevelPriority 根据日志级别创建级别筛选函数
func GetLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	return func(l zapcore.Level) bool {
		return l == level
	}
}

// GetZapEncodeLevel 获取zap日志级别编码器
func GetZapEncodeLevel(encodeLevel string) zapcore.LevelEncoder {
	switch encodeLevel {
	case "LowercaseLevelEncoder":
		return zapcore.LowercaseLevelEncoder
	case "LowercaseColorLevelEncoder":
		return zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder":
		return zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder":
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

// TransportLevel 字符串转日志级别
func TransportLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
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
