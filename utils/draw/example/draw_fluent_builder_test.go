package example

import (
	"fmt"
	"github.com/jiushengTech/common/utils/draw/shape/base"
	"testing"

	"github.com/jiushengTech/common/utils/draw"
)

// TestFluentBuilder 测试改进的API风格
func TestFluentBuilder(t *testing.T) {
	// 测试图片URL
	imageURL := "http://suuqjbby1.hn-bkt.clouddn.com/tryon/origin/uploads/2025-04-22/20250420214131.jpg"

	// 示例1: 使用更简洁的API
	processor := draw.NewImageProcessor(imageURL,
		draw.WithOutputName("fluent_builder_test.png"),
		draw.WithOutputDir("builder_results"),
	)

	// 添加一个圆形
	circle := draw.NewCircle(
		&draw.Point{X: 300, Y: 300},
		50,
		draw.WithColor(base.ColorRed),
		draw.WithFill(true),
	)
	processor.AddShape(circle)

	// 添加一个矩形
	rect := draw.NewRectangle(
		&draw.Point{X: 200, Y: 200},
		&draw.Point{X: 400, Y: 400},
		draw.WithColor(base.ColorBlue),
		draw.WithLineWidth(3),
	)
	processor.AddShape(rect)

	// 处理图像
	outputPath, err := draw.ProcessImage(processor)
	if err != nil {
		t.Errorf("处理图像失败: %v", err)
	} else {
		fmt.Println("改进API - 图片已保存，绝对路径为:", outputPath)
	}

	// 示例2: 使用更组织化的方式
	compositePath, err := processOrganizedWay(imageURL)
	if err != nil {
		t.Errorf("处理图像失败: %v", err)
	} else {
		fmt.Println("组织化API - 图片已保存，绝对路径为:", compositePath)
	}
}

// processOrganizedWay 使用更组织化的方式处理图像
func processOrganizedWay(imageURL string) (string, error) {
	// 创建处理器
	processor := draw.NewImageProcessor(imageURL,
		draw.WithOutputDir("builder_complex_results"),
		draw.WithTimeBasedName(),
		draw.WithOutputFormat(draw.FormatPNG),
	)

	// 创建一组相关的形状
	shapes := []draw.Shape{
		// 创建一个圆形
		draw.NewCircle(
			&draw.Point{X: 400, Y: 300},
			100,
			draw.WithColor(base.ColorRed),
			draw.WithLineWidth(5),
		),
		// 创建一个矩形
		draw.NewRectangle(
			&draw.Point{X: 350, Y: 250},
			&draw.Point{X: 450, Y: 350},
			draw.WithColor(base.ColorBlue),
			draw.WithFill(true),
		),
		// 创建一个线条
		draw.NewLine(
			draw.VerticalLine,
			[]*draw.Point{
				{X: 200, Y: 200},
				{X: 200, Y: 400},
			},
			nil,
			draw.WithColor(base.ColorGreen),
			draw.WithLineWidth(2),
		),
	}

	// 添加所有形状
	for _, shape := range shapes {
		processor.AddShape(shape)
	}

	// 处理图像
	return draw.ProcessImage(processor)
}
