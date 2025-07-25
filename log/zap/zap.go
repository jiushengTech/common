package zap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jiushengTech/common/log/zap/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 默认配置常量
const (
	DefaultLogDir      = "logs"
	DefaultLogLevel    = "info"
	DefaultLogFormat   = "console"
	DefaultMaxSize     = 10 // MB
	DefaultMaxBackups  = 10
	DefaultMaxAge      = 30 // 天
	DefaultEncodeLevel = "LowercaseColorLevelEncoder"
	DefaultStackKey    = "stack"
)

// NewZapLogger 创建并返回一个 zapLogger 实例
func NewZapLogger(c *conf.ZapConf) *zap.Logger {
	// 配置验证
	if c == nil {
		panic("ZapConf 配置不能为 nil")
	}
	if c.Director == "" {
		c.Director = DefaultLogDir
	}
	if c.Level == "" {
		c.Level = DefaultLogLevel
	}
	if c.Format == "" {
		c.Format = DefaultLogFormat
	}
	if c.MaxSize <= 0 {
		c.MaxSize = DefaultMaxSize
	}
	if c.MaxBackups <= 0 {
		c.MaxBackups = DefaultMaxBackups
	}

	cores := getZapCores(c)
	options := []zap.Option{
		zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCallerSkip(int(c.AddCallerSkip)),
	}

	if c.AddCaller {
		options = append(options, zap.AddCaller())
	}
	if c.Model == "dev" {
		options = append(options, zap.Development())
	}
	// 创建并返回 zap.Logger
	return zap.New(zapcore.NewTee(cores...), options...)
}

func DefaultZapLogger() *zap.Logger {
	return NewZapLogger(&conf.ZapConf{
		Model:         "dev",                        // 开发模式配置
		Level:         "debug",                      // 日志级别设置为 debug（捕获 debug、info、warn、error 等）
		Format:        "console",                    // 日志输出格式（console 或 JSON）
		Director:      "logs",                       // 日志文件存储目录
		EncodeLevel:   "LowercaseColorLevelEncoder", // 使用彩色小写级别名称在日志中
		StacktraceKey: "stack",                      // 堆栈跟踪信息的 JSON 键名
		MaxAge:        0,                            // 保留旧日志文件的最大天数（0 表示无限制）
		AddCaller:     true,                         // 显示日志打印所在的行号
		AddCallerSkip: 0,                            // 跳过调用栈的行数
		LogInConsole:  true,                         // 是否在控制台输出日志
		MaxSize:       10,                           // 每个日志文件的最大大小（单位：MB）
		Compress:      false,                        // 是否压缩/归档旧日志文件
		MaxBackups:    10,                           // 保留的旧日志文件的最大数量
		TimeRotation:  RotateHourly,                 // 时间轮转类型: "0:minute", "1:hour" 或 "2:day"
	})
}

// GetEncoder 获取编码器
func GetEncoder(c *conf.ZapConf, forConsole bool) zapcore.Encoder {
	encoderConfig := GetEncoderConfig(c, forConsole)

	if c.Format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	switch c.Format {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// GetEncoderConfig 获取编码器配置
func GetEncoderConfig(c *conf.ZapConf, forConsole bool) zapcore.EncoderConfig {
	encConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller", // 确保记录调用者信息
		StacktraceKey:  c.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 使用FullCallerEncoder，显示完整的文件和行号
		EncodeTime:     CustomTimeEncoder,
	}

	if forConsole {
		// 控制台输出支持彩色编码器
		encConfig.EncodeLevel = getZapEncodeLevel(c.EncodeLevel)
	} else {
		// 文件输出排除彩色编码器
		encConfig.EncodeLevel = getNonColorEncoder(c.EncodeLevel)
	}

	return encConfig
}

// getNonColorEncoder 获取非彩色级别编码器
func getNonColorEncoder(encodeLevel string) zapcore.LevelEncoder {
	// 无论输入是彩色还是非彩色编码器选项，都返回对应的非彩色版本
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
func CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format(time.DateTime))
}

// getZapCores 获取所有zap核心组件
func getZapCores(c *conf.ZapConf) []zapcore.Core {
	cores := make([]zapcore.Core, 0, 2)
	minLevel := TransportLevel(c.Level)

	// 为每个级别创建单独的Core，提高性能
	levels := []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		zapcore.DPanicLevel,
		zapcore.PanicLevel,
		zapcore.FatalLevel,
	}

	for _, level := range levels {
		if level >= minLevel {
			// 为每个级别创建单独的文件Core
			fileCore := zapcore.NewCore(
				GetEncoder(c, false),
				GetWriteSyncer(c, level.String()),
				GetLevelPriority(level),
			)
			cores = append(cores, fileCore)
		}
	}

	// 如果需要控制台输出，添加控制台Core
	if c.LogInConsole {
		levelFunc := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= minLevel
		})
		consoleCore := createConsoleCore(c, levelFunc)
		cores = append(cores, consoleCore)
	}

	return cores
}

// createConsoleCore 创建控制台输出的Core
func createConsoleCore(c *conf.ZapConf, levelFunc zap.LevelEnablerFunc) zapcore.Core {
	isColorful := true
	if strings.ToLower(c.Format) == "json" {
		isColorful = false // JSON 格式不使用彩色编码器，避免乱码
	}
	return zapcore.NewCore(
		GetEncoder(c, isColorful),
		zapcore.AddSync(os.Stdout),
		levelFunc,
	)
}

// GetWriteSyncer 创建文件日志写入器，支持按照大小和时间切割
func GetWriteSyncer(c *conf.ZapConf, level string) zapcore.WriteSyncer {
	// 创建日志目录
	logPath := filepath.Join(c.Director, time.Now().Format("2006-01"))
	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		panic(fmt.Sprintf("创建日志目录失败: %v", err))
	}

	// 创建级别目录
	levelDir := filepath.Join(logPath, level)
	if err := os.MkdirAll(levelDir, os.ModePerm); err != nil {
		panic(fmt.Sprintf("创建日志级别目录失败: %v", err))
	}

	writer := NewTimeRotationWriter(c, level, levelDir)
	return zapcore.AddSync(writer)
}

// GetLevelPriority 根据日志级别创建级别筛选函数
func GetLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	return func(l zapcore.Level) bool {
		return l == level
	}
}

// getZapEncodeLevel 获取zap日志级别编码器
func getZapEncodeLevel(encodeLevel string) zapcore.LevelEncoder {
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
