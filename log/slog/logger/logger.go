// Package logger 提供基于Go的slog构建的灵活日志解决方案
// 支持按时间（分钟/小时/天）和大小进行日志轮转
//
// 它提供简洁的API和强大功能：
// - 多种轮转策略（分钟、小时、天、大小）
// - 线程安全的日志记录
// - 支持文件和标准输出
// - JSON和文本格式选项
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// RotateMode 定义日志轮转模式类型
type RotateMode int

const (
	// RotateDaily 每天轮转一次日志
	RotateDaily RotateMode = iota
	// RotateHourly 每小时轮转一次日志
	RotateHourly
	// RotateMinutely 每分钟轮转一次日志
	RotateMinutely
)

// String 返回轮转模式的字符串表示
func (m RotateMode) String() string {
	switch m {
	case RotateDaily:
		return "daily"
	case RotateHourly:
		return "hourly"
	case RotateMinutely:
		return "minutely"
	default:
		return "unknown"
	}
}

// Logger 是日志记录器的主要结构体，封装了slog.Logger
// 提供自动轮转功能和多种输出选项
type Logger struct {
	log         *slog.Logger       // 底层slog记录器
	opts        *Options           // 配置选项
	mu          sync.Mutex         // 互斥锁，保证并发安全
	fileDate    string             // 当前日志文件的日期
	fileHour    int                // 当前日志文件的小时值
	fileMinute  int                // 当前日志文件的分钟值
	rotateMode  RotateMode         // 当前轮转模式
	lumberjack  *lumberjack.Logger // 用于文件轮转的lumberjack实例
	currentSize int64              // 当前日志文件的估计大小
}

// Options 定义日志记录器的配置选项
type Options struct {
	// 基本选项
	Level       slog.Level                                      // 日志级别
	Format      string                                          // 日志格式: "json" 或 "text"
	AddSource   bool                                            // 是否添加源代码位置信息
	Writers     []io.Writer                                     // 自定义写入器列表
	ReplaceAttr func(groups []string, attr slog.Attr) slog.Attr // 自定义属性替换函数

	// 文件选项
	LogDir       string // 日志目录
	FilePrefix   string // 日志文件前缀
	EnableFile   bool   // 是否启用文件日志
	EnableStdout bool   // 是否启用标准输出

	// 轮转选项
	MaxSize    int        // 单个日志文件最大大小（MB）
	MaxBackups int        // 最大保留的历史日志文件数
	MaxAge     int        // 日志文件最大保留天数
	Compress   bool       // 是否压缩历史日志
	RotateMode RotateMode // 日志轮转模式
}

// DefaultOptions 返回默认的日志配置选项
// 默认为：文本格式、INFO级别、同时输出到文件和控制台、按小时轮转
func DefaultOptions() *Options {
	return &Options{
		// 基本配置
		Level:     slog.LevelInfo,
		Format:    "text",
		AddSource: true,
		Writers:   []io.Writer{os.Stdout},

		// 文件配置
		LogDir:       "./logs",
		FilePrefix:   "app",
		EnableFile:   true,
		EnableStdout: true,

		// 轮转配置
		MaxSize:    10, // 10MB
		MaxBackups: 5,  // 保留5个备份
		MaxAge:     7,  // 保留7天
		Compress:   false,
		RotateMode: RotateHourly, // 默认按小时轮转
	}
}

// New 创建并返回一个新的Logger实例
// 如果opts为nil，则使用默认选项
func New(opts *Options) (*Logger, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	logger := &Logger{
		opts:       opts,
		rotateMode: opts.RotateMode,
	}

	// 确保日志目录存在
	if opts.EnableFile {
		if err := os.MkdirAll(opts.LogDir, 0755); err != nil {
			return nil, fmt.Errorf("创建日志目录失败: %w", err)
		}
	}

	// 初始化日志输出
	if err := logger.initLogWriter(); err != nil {
		return nil, fmt.Errorf("初始化日志写入器失败: %w", err)
	}

	return logger, nil
}

// getCurrentLogFileName 根据轮转模式生成当前的日志文件名
func (l *Logger) getCurrentLogFileName() string {
	now := time.Now()
	currentDate := now.Format(time.DateOnly)

	switch l.rotateMode {
	case RotateMinutely:
		// 按分钟命名: prefix-YYYY-MM-DD-HH-MM.log
		return fmt.Sprintf("%s-%s-%02d-%02d.log", l.opts.FilePrefix, currentDate, now.Hour(), now.Minute())
	case RotateHourly:
		// 按小时命名: prefix-YYYY-MM-DD-HH.log
		return fmt.Sprintf("%s-%s-%02d.log", l.opts.FilePrefix, currentDate, now.Hour())
	default:
		// 按天命名: prefix-YYYY-MM-DD.log
		return fmt.Sprintf("%s-%s.log", l.opts.FilePrefix, currentDate)
	}
}

// initLogWriter 初始化日志写入器
// 根据配置设置输出目标和格式
func (l *Logger) initLogWriter() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	currentDate := now.Format(time.DateOnly)
	currentHour := now.Hour()
	currentMinute := now.Minute()

	var writers []io.Writer

	// 配置文件输出
	if l.opts.EnableFile {
		logFileName := l.getCurrentLogFileName()
		logFilePath := filepath.Join(l.opts.LogDir, logFileName)

		// 配置lumberjack进行文件轮转
		l.lumberjack = &lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    l.opts.MaxSize,
			MaxBackups: l.opts.MaxBackups,
			MaxAge:     l.opts.MaxAge,
			Compress:   l.opts.Compress,
			LocalTime:  true, // 使用本地时间
		}

		writers = append(writers, l.lumberjack)
		l.fileDate = currentDate
		l.fileHour = currentHour
		l.fileMinute = currentMinute

		// 获取当前文件大小
		if fileInfo, err := os.Stat(logFilePath); err == nil {
			l.currentSize = fileInfo.Size()
		} else {
			l.currentSize = 0
		}
	}

	// 配置标准输出
	if l.opts.EnableStdout {
		writers = append(writers, os.Stdout)
	}

	// 创建多路输出
	var writer io.Writer
	if len(writers) > 1 {
		writer = io.MultiWriter(writers...)
	} else if len(writers) == 1 {
		writer = writers[0]
	} else {
		// 没有配置输出目标时，默认使用标准输出
		writer = os.Stdout
	}

	// 根据配置选择日志格式
	var handler slog.Handler
	switch l.opts.Format {
	case "json":
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level:       l.opts.Level,
			AddSource:   l.opts.AddSource,
			ReplaceAttr: l.opts.ReplaceAttr,
		})
	default:
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{
			Level:       l.opts.Level,
			AddSource:   l.opts.AddSource,
			ReplaceAttr: l.opts.ReplaceAttr,
		})
	}

	// 创建新的slog.Logger
	l.log = slog.New(handler)
	return nil
}

// checkRotate 检查是否需要轮转日志
// 同时检查时间条件和文件大小条件
func (l *Logger) checkRotate() error {
	now := time.Now()
	currentDate := now.Format(time.DateOnly)
	currentHour := now.Hour()
	currentMinute := now.Minute()

	// 检查是否需要按时间轮转
	needsRotate := false
	switch l.rotateMode {
	case RotateMinutely:
		if currentDate != l.fileDate || currentHour != l.fileHour || currentMinute != l.fileMinute {
			needsRotate = true
		}
	case RotateHourly:
		if currentDate != l.fileDate || currentHour != l.fileHour {
			needsRotate = true
		}
	case RotateDaily:
		if currentDate != l.fileDate {
			needsRotate = true
		}
	}

	// 如果不需要按时间轮转，则检查是否需要按大小轮转
	if !needsRotate && l.opts.EnableFile && l.lumberjack != nil {
		// 每次写入估算增加1KB，减少实际检查文件大小的频率
		l.currentSize += 1024 // 假设每次日志条目约1KB

		// 当估算大小接近最大值时，检查实际文件大小
		maxSizeBytes := int64(l.opts.MaxSize) * 1024 * 1024
		if l.currentSize >= maxSizeBytes-1024*10 { // 当接近阈值时检查
			logFilePath := filepath.Join(l.opts.LogDir, l.getCurrentLogFileName())
			if fileInfo, err := os.Stat(logFilePath); err == nil {
				l.currentSize = fileInfo.Size()
				if l.currentSize >= maxSizeBytes {
					needsRotate = true
				}
			}
		}
	}

	// 如果需要轮转，重新初始化日志写入器
	if needsRotate {
		return l.initLogWriter()
	}

	return nil
}

// Write 实现io.Writer接口
// 允许直接将Logger用作io.Writer
func (l *Logger) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.lumberjack != nil {
		n, err = l.lumberjack.Write(p)
		if err == nil {
			l.currentSize += int64(n)
		}
		return n, err
	}

	return 0, fmt.Errorf("lumberjack未初始化")
}

// Debug 记录Debug级别的日志
func (l *Logger) Debug(msg string, args ...any) {
	if err := l.checkRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "日志轮转失败: %v\n", err)
	}
	l.log.Debug(msg, args...)
}

// Info 记录Info级别的日志
func (l *Logger) Info(msg string, args ...any) {
	if err := l.checkRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "日志轮转失败: %v\n", err)
	}
	l.log.Info(msg, args...)
}

// Warn 记录Warn级别的日志
func (l *Logger) Warn(msg string, args ...any) {
	if err := l.checkRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "日志轮转失败: %v\n", err)
	}
	l.log.Warn(msg, args...)
}

// Error 记录Error级别的日志
func (l *Logger) Error(msg string, args ...any) {
	if err := l.checkRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "日志轮转失败: %v\n", err)
	}
	l.log.Error(msg, args...)
}

// DebugContext 记录带上下文的Debug级别日志
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	if err := l.checkRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "日志轮转失败: %v\n", err)
	}
	l.log.DebugContext(ctx, msg, args...)
}

// InfoContext 记录带上下文的Info级别日志
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	if err := l.checkRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "日志轮转失败: %v\n", err)
	}
	l.log.InfoContext(ctx, msg, args...)
}

// WarnContext 记录带上下文的Warn级别日志
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	if err := l.checkRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "日志轮转失败: %v\n", err)
	}
	l.log.WarnContext(ctx, msg, args...)
}

// ErrorContext 记录带上下文的Error级别日志
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	if err := l.checkRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "日志轮转失败: %v\n", err)
	}
	l.log.ErrorContext(ctx, msg, args...)
}

// Close 关闭日志记录器
// 虽然lumberjack不需要显式关闭，但提供此方法以兼容io.Closer接口
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return nil
}

// SetRotateMode 动态设置日志轮转模式
func (l *Logger) SetRotateMode(mode RotateMode) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.rotateMode != mode {
		l.rotateMode = mode
		l.opts.RotateMode = mode
		// 更改模式后重新初始化日志写入器
		_ = l.initLogWriter()
	}
}

// SetMaxSize 动态设置日志文件最大大小（MB）
func (l *Logger) SetMaxSize(maxSize int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.opts.MaxSize != maxSize {
		l.opts.MaxSize = maxSize
		if l.lumberjack != nil {
			l.lumberjack.MaxSize = maxSize
		}
	}
}

// GetSlogLogger 获取底层的slog.Logger实例
// 允许用户直接使用slog的高级功能
func (l *Logger) GetSlogLogger() *slog.Logger {
	return l.log
}

// SetLevel 动态设置日志级别
func (l *Logger) SetLevel(level slog.Level) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.opts.Level = level
	return l.initLogWriter()
}

// GetOptions 获取当前配置选项的副本
func (l *Logger) GetOptions() Options {
	l.mu.Lock()
	defer l.mu.Unlock()
	return *l.opts
}
