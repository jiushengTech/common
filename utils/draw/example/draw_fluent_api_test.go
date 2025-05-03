package example

import (
	"fmt"
	"github.com/jiushengTech/common/utils/draw/shape/base"
	"testing"

	"github.com/jiushengTech/common/utils/draw"
)

// TestDrawWithFluentAPI 测试流畅的API设计
func TestDrawWithFluentAPI(t *testing.T) {
	// 测试图片URL
	imageURL := "http://suuqjbby1.hn-bkt.clouddn.com/tryon/origin/uploads/2025-04-22/20250420214131.jpg"

	// 使用流畅API处理图像
	outputPath, err := processWithFluentAPI(imageURL)
	if err != nil {
		t.Errorf("处理图像失败: %v", err)
	} else {
		fmt.Println("流畅API - 图片已保存，绝对路径为:", outputPath)
	}

	// 测试更复杂的图形组合和变换
	compositePath, err := processWithCompositeShapes(imageURL)
	if err != nil {
		t.Errorf("处理图像失败: %v", err)
	} else {
		fmt.Println("复杂组合 - 图片已保存，绝对路径为:", compositePath)
	}
}

// processWithFluentAPI 使用链式调用风格处理图像
func processWithFluentAPI(imageURL string) (string, error) {
	// 创建一个图像处理器构造器
	processor := draw.NewImageProcessor(
		imageURL,
		draw.WithOutputDir("fluent_api_results"),
		draw.WithTimeBasedName(),
		draw.WithOutputFormat(draw.FormatPNG),
	)

	// 添加一个圆形
	circle := draw.NewCircle(
		&draw.Point{X: 400, Y: 300},
		100,
		draw.WithColor(base.ColorRed),
		draw.WithLineWidth(5),
	)
	processor.AddShape(circle)

	// 添加一个矩形
	rect := draw.NewRectangle(
		&draw.Point{X: 200, Y: 200},
		&draw.Point{X: 600, Y: 400},
		draw.WithColor(base.ColorBlue),
		draw.WithLineWidth(3),
		draw.WithFill(false),
	)
	processor.AddShape(rect)

	// 处理图像并获取路径
	return draw.ProcessImage(processor)
}

// processWithCompositeShapes 使用图形组合和更高级的功能
func processWithCompositeShapes(imageURL string) (string, error) {
	// 创建处理器
	processor := draw.NewImageProcessor(
		imageURL,
		draw.WithOutputDir("composite_results"),
		draw.WithTimeBasedName(),
		draw.WithJpegQuality(95),
		draw.WithOutputFormat(draw.FormatJPEG),
	)

	// 创建形状组
	group := draw.NewShapeGroup("highlights")

	// 添加多个圆形到组
	for i := 0; i < 5; i++ {
		circle := draw.NewCircle(
			&draw.Point{X: 300 + float64(i*60), Y: 200},
			20,
			draw.WithColor(base.ColorRed),
			draw.WithFill(true),
		)
		group.AddShape(circle)
	}

	// 添加矩形到组
	rect := draw.NewRectangle(
		&draw.Point{X: 250, Y: 250},
		&draw.Point{X: 550, Y: 350},
		draw.WithColor(base.ColorBlue),
		draw.WithLineWidth(2),
	)
	group.AddShape(rect)

	// 将组中的所有形状添加到处理器
	for _, shape := range group.Shapes {
		processor.AddShape(shape)
	}

	// 处理图像并获取路径
	return draw.ProcessImage(processor)
}
