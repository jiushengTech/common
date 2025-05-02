package main

import (
	"fmt"
	"testing"

	"github.com/jiushengTech/common/utils/draw"
)

// TestAbsolutePathWithShapes 测试绘制图形并返回绝对路径
func TestAbsolutePathWithShapes(t *testing.T) {
	// 测试图片URL (使用一个示例图片)
	imageURL := "http://suuqjbby1.hn-bkt.clouddn.com/tryon/origin/uploads/2025-04-22/20250420214131.jpg"

	// 创建一个圆形
	circle := draw.NewCircle(
		draw.Point{X: 400, Y: 300},
		100,
		draw.WithColor(draw.ColorRed),
		draw.WithLineWidth(5),
	)

	// 创建一个矩形
	rect := draw.NewRectangle(
		draw.Point{X: 200, Y: 200},
		draw.Point{X: 600, Y: 400},
		draw.WithColor(draw.ColorBlue),
		draw.WithLineWidth(3),
		draw.WithFill(false),
	)

	// 创建图像处理器
	processor := draw.NewImageProcessor(
		imageURL,
		draw.WithTimeBasedName(),
		draw.WithOutputDir("group_test_results"),
		draw.WithShape(circle),
		draw.WithShape(rect),
	)

	// 方法1：处理图像并直接获取绝对路径
	absPath, err := draw.ProcessImage(processor)
	if err != nil {
		t.Errorf("处理图像失败: %v", err)
	} else {
		fmt.Println("方法1 - 图片已保存，绝对路径为:", absPath)
	}

	// 创建另一个处理器用于方法2演示
	processor2 := draw.NewImageProcessor(
		imageURL,
		draw.WithTimeBasedName(),
		draw.WithOutputDir("group_test_results"),
		draw.WithShape(circle),
	)

	// 方法2：先处理图像，然后获取绝对路径
	absPath2, err := processor2.Process()
	if err != nil {
		t.Errorf("处理图像失败: %v", err)
	} else {
		fmt.Println("方法2 - 图片已保存，绝对路径为:", absPath2)

		// 也可以使用 GetAbsoluteOutputPath 获取路径（在某些情况下更有用）
		pathFromGetter, err := draw.GetAbsoluteOutputPath(processor2)
		if err != nil {
			t.Errorf("获取绝对路径失败: %v", err)
		} else {
			fmt.Println("方法2 (通过GetAbsoluteOutputPath) - 图片绝对路径为:", pathFromGetter)
		}
	}

	// 清理临时文件
	draw.CleanupAllTempFiles()
}
