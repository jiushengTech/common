package example

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jiushengTech/common/utils/draw"
)

// TestDefaultOutputName 测试默认输出文件名是否基于当前时间
func TestDefaultOutputName(t *testing.T) {
	// 测试图片URL
	imageURL := "https://picsum.photos/300/200"

	// 创建一个简单的图形
	circle := draw.NewCircle(
		&draw.Point{X: 150, Y: 100},
		50,
		draw.WithColor(draw.ColorBlue),
		draw.WithFill(true),
	)

	// 创建图像处理器，不指定文件名
	processor := draw.NewImageProcessor(
		imageURL,
		draw.WithOutputDir("default_name_test"),
		draw.WithShape(circle),
	)

	// 获取当前时间，用于验证文件名格式
	now := time.Now()
	expectedPrefix := fmt.Sprintf("%d%02d%02d_",
		now.Year(), now.Month(), now.Day())

	// 检查默认文件名是否基于当前时间
	if !strings.HasPrefix(processor.Output, expectedPrefix) {
		t.Errorf("期望文件名以 %s 开头，但实际是 %s", expectedPrefix, processor.Output)
	} else {
		fmt.Printf("默认文件名格式正确: %s\n", processor.Output)
	}

	// 处理图像
	absPath, err := processor.Process()
	if err != nil {
		t.Errorf("处理图像失败: %v", err)
	} else {
		fmt.Printf("图片已保存，绝对路径为: %s\n", absPath)
	}

	// 清理临时文件
	draw.CleanupAllTempFiles()
}
