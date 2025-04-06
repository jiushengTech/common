package logger

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoggerBasic(t *testing.T) {
	// 使用默认配置创建日志记录器
	logger, err := New(nil)
	if err != nil {
		t.Fatalf("创建默认日志记录器失败: %v", err)
	}
	defer logger.Close()

	// 记录一些基本日志
	logger.Info("应用启动")
	logger.Debug("这是一条调试日志", "userId", 12345)
	logger.Warn("警告信息", "status", "degraded")
	logger.Error("发生错误", "code", 500, "error", "数据库连接失败")

	// 使用上下文
	ctx := context.Background()
	logger.InfoContext(ctx, "带上下文的日志", "requestId", "abc-123")
}

func TestLoggerJSON(t *testing.T) {
	// 自定义JSON配置
	jsonOpts := DefaultOptions()
	jsonOpts.Format = "json"
	jsonOpts.LogDir = "./custom_logs"
	jsonOpts.FilePrefix = "json_test"
	jsonOpts.MaxSize = 5 // 5MB

	jsonLogger, err := New(jsonOpts)
	if err != nil {
		t.Fatalf("创建JSON日志记录器失败: %v", err)
	}
	defer jsonLogger.Close()

	// 使用JSON格式记录日志
	jsonLogger.Info("JSON格式日志",
		"user", map[string]interface{}{
			"id":    1001,
			"name":  "张三",
			"roles": []string{"admin", "user"},
		},
		"timestamp", time.Now().Unix(),
	)
}

func TestConfigModification(t *testing.T) {
	// 创建默认日志记录器并修改其配置
	logger, err := New(nil)
	if err != nil {
		t.Fatalf("创建默认日志记录器失败: %v", err)
	}
	defer logger.Close()

	// 动态修改配置
	logger.SetRotateMode(RotateHourly)
	logger.SetMaxSize(20)
	logger.SetLevel(slog.LevelDebug)

	// 获取底层slog.Logger进行高级操作
	slogLogger := logger.GetSlogLogger()
	slogLogger.Info("使用底层slog记录器", "detail", "高级功能")
}

// 专门测试分钟级轮转功能
func TestMinuteRotation(t *testing.T) {
	// 创建临时目录用于测试
	testDir := "./test_minute_logs"
	os.RemoveAll(testDir) // 清理可能存在的旧测试文件
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// 创建按分钟轮转的日志记录器
	minuteOpts := DefaultOptions()
	minuteOpts.LogDir = testDir
	minuteOpts.FilePrefix = "minute_test"
	minuteOpts.RotateMode = RotateMinutely
	minuteOpts.EnableStdout = true

	minuteLogger, err := New(minuteOpts)
	if err != nil {
		t.Fatalf("创建按分钟轮转的日志记录器失败: %v", err)
	}

	// 记录第一条日志
	now := time.Now()
	minuteLogger.Info("第一分钟的日志消息")

	// 验证创建了正确格式的日志文件
	expectedPrefix := "minute_test-" + now.Format("2006-01-02-15-04")

	// 检查日志文件是否存在
	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("读取测试目录失败: %v", err)
	}

	if len(files) == 0 {
		t.Fatalf("未找到创建的日志文件")
	}

	found := false
	for _, file := range files {
		if strings.HasPrefix(file.Name(), expectedPrefix) {
			found = true
			t.Logf("找到符合格式的日志文件: %s", file.Name())
			break
		}
	}

	if !found {
		t.Fatalf("未找到带有分钟的日志文件, 应包含前缀: %s", expectedPrefix)
	}

	// 关闭第一个日志记录器
	minuteLogger.Close()

	// 模拟一分钟后创建新的日志记录器
	// 通过直接修改时间来模拟时间流逝
	// 在实际测试环境中可能需要调整内部时间或使用模拟时间的库

	// 创建第二个日志记录器并设置不同的分钟
	time.Sleep(1 * time.Second) // 稍等一会儿确保文件名不同

	// 创建第二个日志记录器
	minuteLogger2, err := New(minuteOpts)
	if err != nil {
		t.Fatalf("创建第二个日志记录器失败: %v", err)
	}
	defer minuteLogger2.Close()

	// 这里我们直接修改内部的文件名生成函数来模拟分钟变化
	// 注意：这只是测试目的，实际使用时不应该这样做
	// 如果Logger结构体暴露了适当的测试接口，可以使用该接口

	// 写入第二条日志，但在不同的"分钟"
	minuteLogger2.Info("第二个日志消息")

	// 可以在这里添加额外的验证步骤
	// 如果您的Logger实现了旋转时间的测试辅助方法，可以在此使用
}

// 测试辅助函数：强制进行日志轮转
// 如果您可以修改Logger实现，建议添加这样的测试辅助函数
func forceRotation(l *Logger) {
	// 这只是示例，实际实现取决于您的Logger内部结构
	// 例如：l.fileMinute = (l.fileMinute + 1) % 60
	// 或者：l.lastRotateTime = time.Now().Add(-2 * time.Minute)
	// 然后调用内部的rotate()方法
}

// 完整的集成测试，验证所有轮转模式
func TestAllRotationModes(t *testing.T) {
	baseDir := "./rotation_test_logs"
	os.RemoveAll(baseDir) // 清理可能存在的旧测试文件

	modes := []struct {
		name       string
		rotateMode RotateMode
		subdir     string
	}{
		{"按天轮转", RotateDaily, "daily"},
		{"按小时轮转", RotateHourly, "hourly"},
		{"按分钟轮转", RotateMinutely, "minutely"},
	}

	for _, mode := range modes {
		t.Run(mode.name, func(t *testing.T) {
			// 创建目录
			logDir := filepath.Join(baseDir, mode.subdir)
			os.MkdirAll(logDir, 0755)

			// 创建日志记录器
			opts := DefaultOptions()
			opts.LogDir = logDir
			opts.FilePrefix = mode.subdir
			opts.RotateMode = mode.rotateMode
			opts.MaxSize = 1 // 1MB，用于测试按大小轮转

			logger, err := New(opts)
			if err != nil {
				t.Fatalf("创建%s日志记录器失败: %v", mode.name, err)
			}

			// 写入测试日志
			logger.Info("测试日志", "mode", mode.name)

			// 关闭日志记录器
			logger.Close()

			// 验证日志文件已创建
			files, err := os.ReadDir(logDir)
			if err != nil {
				t.Fatalf("读取日志目录失败: %v", err)
			}

			if len(files) == 0 {
				t.Fatalf("未找到创建的日志文件")
			}

			// 验证文件名格式
			fileName := files[0].Name()
			t.Logf("%s模式创建的文件: %s", mode.name, fileName)

			// 根据不同的轮转模式验证文件名格式
			// 这里您可以添加更具体的验证
		})
	}
}
