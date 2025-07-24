package zap

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/jiushengTech/common/log/zap/conf"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// RotateMinutely 每分钟轮转一次日志
	RotateMinutely = iota
	// RotateHourly 每小时轮转一次日志
	RotateHourly
	// RotateDaily 每天轮转一次日志
	RotateDaily
)

// TimeRotationInfo 存储日志文件轮转的相关信息
type TimeRotationInfo struct {
	CurrentFileName string    // 当前日志文件名
	NextRotation    time.Time // 下一次轮转时间
	LastAccessed    time.Time // 最后访问时间，用于清理无用的跟踪器
}

// 全局映射表，用于跟踪每个日志文件的轮转信息
var (
	rotationTrackers = make(map[string]*TimeRotationInfo)
	trackerMutex     sync.RWMutex // 保护映射表的并发访问
)

// TimeRotationHook 实现了io.WriteCloser接口，用于处理时间轮转
type TimeRotationHook struct {
	Lumberjack      *lumberjack.Logger
	RotationTracker *TimeRotationInfo
	Config          *conf.ZapConf
	Level           string
	LevelDir        string
}

// Write 写入日志，并检查是否需要轮转
func (t *TimeRotationHook) Write(p []byte) (n int, err error) {
	now := time.Now()

	// 使用读锁检查是否需要轮转
	trackerMutex.RLock()
	needsRotation := now.After(t.RotationTracker.NextRotation)
	trackerMutex.RUnlock()

	if needsRotation {
		// 获取写锁进行轮转操作
		trackerMutex.Lock()
		// 双重检查，避免多个goroutine同时轮转
		if now.After(t.RotationTracker.NextRotation) {
			// 关闭旧文件（重要）
			_ = t.Lumberjack.Close()
			// 生成新文件名和轮转时间
			logFileName, nextRotation := generateFileNameAndRotation(now, int(t.Config.TimeRotation), t.Level)
			logFilePath := filepath.Join(t.LevelDir, logFileName)
			// 创建新的 lumberjack 实例
			t.Lumberjack = &lumberjack.Logger{
				Filename:   logFilePath,
				MaxSize:    int(t.Config.MaxSize),
				MaxAge:     int(t.Config.MaxAge),
				MaxBackups: int(t.Config.MaxBackups),
				LocalTime:  true,
				Compress:   t.Config.Compress,
			}
			// 更新轮转状态
			t.RotationTracker.CurrentFileName = logFileName
			t.RotationTracker.NextRotation = nextRotation
			t.RotationTracker.LastAccessed = now
		}
		trackerMutex.Unlock()
	} else {
		// 更新最后访问时间
		trackerMutex.Lock()
		t.RotationTracker.LastAccessed = now
		trackerMutex.Unlock()
	}

	// 正常写入日志
	return t.Lumberjack.Write(p)
}

// Close 关闭日志文件
func (t *TimeRotationHook) Close() error {
	return t.Lumberjack.Close()
}

// Sync 同步日志到磁盘
func (t *TimeRotationHook) Sync() error {
	return nil // lumberjack已经处理了同步
}

// NewTimeRotationWriter 创建一个支持时间轮转的日志写入器
func NewTimeRotationWriter(c *conf.ZapConf, level, levelDir string) *TimeRotationHook {
	now := time.Now()

	// 生成文件名和计算下一次轮转时间
	logFileName, nextRotation := generateFileNameAndRotation(now, int(c.TimeRotation), level)

	logFilePath := filepath.Join(levelDir, logFileName)

	// 更新或创建轮转跟踪信息
	trackerKey := levelDir + "-" + level

	trackerMutex.Lock()
	defer trackerMutex.Unlock()

	// 清理过期的跟踪器（超过24小时未访问）
	cleanupOldTrackers(now)

	if _, exists := rotationTrackers[trackerKey]; !exists || rotationTrackers[trackerKey].CurrentFileName != logFileName {
		rotationTrackers[trackerKey] = &TimeRotationInfo{
			CurrentFileName: logFileName,
			NextRotation:    nextRotation,
			LastAccessed:    now,
		}
	}

	// 配置lumberjack日志切割器
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    int(c.MaxSize),
		MaxAge:     int(c.MaxAge),
		MaxBackups: int(c.MaxBackups),
		LocalTime:  true,
		Compress:   c.Compress,
	}

	// 返回时间轮转钩子
	return &TimeRotationHook{
		Lumberjack:      lumberjackLogger,
		RotationTracker: rotationTrackers[trackerKey],
		Config:          c,
		Level:           level,
		LevelDir:        levelDir,
	}
}

// cleanupOldTrackers 清理超过24小时未访问的跟踪器
func cleanupOldTrackers(now time.Time) {
	cutoff := now.Add(-24 * time.Hour)
	for key, tracker := range rotationTrackers {
		if tracker.LastAccessed.Before(cutoff) {
			delete(rotationTrackers, key)
		}
	}
}

// generateFileNameAndRotation 根据轮转类型生成文件名和下一次轮转时间
func generateFileNameAndRotation(now time.Time, rotationType int, level string) (string, time.Time) {
	var logFileName string
	var nextRotation time.Time
	switch rotationType {
	case RotateMinutely:
		// 分钟级别轮转 - 文件名格式: 2006-01-02-15-04-level.log
		logFileName = fmt.Sprintf("%s-%02d-%02d-%s.log", now.Format(time.DateOnly), now.Hour(), now.Minute(), level)
		nextRotation = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+1, 0, 0, now.Location())
	case RotateHourly:
		// 小时级别轮转 - 文件名格式: 2006-01-02-15-level.log
		logFileName = fmt.Sprintf("%s-%02d-%s.log", now.Format(time.DateOnly), now.Hour(), level)
		nextRotation = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
	case RotateDaily:
		// 天级别轮转 - 文件名格式: 2006-01-02-level.log
		logFileName = fmt.Sprintf("%s-%s.log", now.Format(time.DateOnly), level)
		nextRotation = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	default:
		// 默认按照天切割
		logFileName = fmt.Sprintf("%s-%s.log", now.Format(time.DateOnly), level)
		nextRotation = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	}

	return logFileName, nextRotation
}
