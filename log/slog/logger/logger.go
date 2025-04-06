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

// Level 定义日志级别
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// Logger 结构体，封装 slog.Logger
type Logger struct {
	log      *slog.Logger
	opts     *Options
	mu       sync.Mutex
	fileDate string
}

// Options 日志配置选项
type Options struct {
	Level        slog.Level
	Format       string // "json" or "text"
	AddSource    bool
	Writers      []io.Writer
	ReplaceAttr  func(groups []string, attr slog.Attr) slog.Attr
	LogDir       string // 日志目录
	FilePrefix   string // 日志文件前缀
	EnableFile   bool   // 是否启用文件日志
	EnableStdout bool   // 是否启用标准输出
	MaxSize      int    // 单个日志文件最大大小（MB）
	MaxBackups   int    // 最多保留多少个历史日志
	MaxAge       int    // 日志最多保存多少天
	Compress     bool   // 是否压缩历史日志
}

// DefaultOptions 返回默认日志配置
func DefaultOptions() *Options {
	return &Options{
		Level:        slog.LevelInfo,
		Format:       "text",
		AddSource:    true,
		Writers:      []io.Writer{os.Stdout},
		LogDir:       "./logs",
		FilePrefix:   "app",
		EnableFile:   true,
		EnableStdout: true,
		MaxSize:      10, // 10MB
		MaxBackups:   5,
		MaxAge:       7,
		Compress:     false,
	}
}

// NewSlogLogger 创建 Logger
func NewSlogLogger(opts *Options) (*Logger, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	logger := &Logger{
		opts: opts,
	}

	// 确保日志目录存在
	if opts.EnableFile {
		if err := os.MkdirAll(opts.LogDir, 0755); err != nil {
			return nil, fmt.Errorf("创建日志目录失败: %w", err)
		}
	}

	// 初始化日志输出
	if err := logger.initLogWriter(); err != nil {
		return nil, err
	}

	return logger, nil
}

// initLogWriter 初始化日志写入
func (l *Logger) initLogWriter() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	currentDate := time.Now().Format("2006-01-02")
	if currentDate == l.fileDate {
		return nil
	}

	var writers []io.Writer

	// 配置日志文件
	if l.opts.EnableFile {
		logFilePath := filepath.Join(l.opts.LogDir, fmt.Sprintf("%s-%s.log", l.opts.FilePrefix, currentDate))
		logWriter := &lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    l.opts.MaxSize,
			MaxBackups: l.opts.MaxBackups,
			MaxAge:     l.opts.MaxAge,
			Compress:   l.opts.Compress,
			LocalTime:  true, // 让日志文件使用本地时间
		}
		writers = append(writers, logWriter)
		l.fileDate = currentDate
	}

	// 配置标准输出
	if l.opts.EnableStdout {
		writers = append(writers, os.Stdout)
	}

	// 创建 MultiWriter
	var writer io.Writer
	if len(writers) > 1 {
		writer = io.MultiWriter(writers...)
	} else if len(writers) == 1 {
		writer = writers[0]
	} else {
		writer = os.Stdout
	}

	// 选择日志格式
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

	l.log = slog.New(handler)
	return nil
}

// checkRotate 检查是否需要轮转日志
func (l *Logger) checkRotate() error {
	currentDate := time.Now().Format("2006-01-02")
	if currentDate != l.fileDate {
		return l.initLogWriter()
	}
	return nil
}

// Debug 记录 Debug 级别日志
func (l *Logger) Debug(msg string, args ...any) {
	if err := l.checkRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "日志轮转失败: %v\n", err)
	}
	l.log.Debug(msg, args...)
}

// Info 记录 Info 级别日志
func (l *Logger) Info(msg string, args ...any) {
	if err := l.checkRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "日志轮转失败: %v\n", err)
	}
	l.log.Info(msg, args...)
}

// Warn 记录 Warn 级别日志
func (l *Logger) Warn(msg string, args ...any) {
	if err := l.checkRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "日志轮转失败: %v\n", err)
	}
	l.log.Warn(msg, args...)
}

// Error 记录 Error 级别日志
func (l *Logger) Error(msg string, args ...any) {
	if err := l.checkRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "日志轮转失败: %v\n", err)
	}
	l.log.Error(msg, args...)
}

// DebugContext 记录 Debug 级别日志（带 Context）
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	if err := l.checkRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "日志轮转失败: %v\n", err)
	}
	l.log.DebugContext(ctx, msg, args...)
}

// InfoContext 记录 Info 级别日志（带 Context）
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	if err := l.checkRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "日志轮转失败: %v\n", err)
	}
	l.log.InfoContext(ctx, msg, args...)
}

// Close 关闭日志文件（如果使用了文件）
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return nil // lumberjack 无需显式关闭
}
