package main

import (
	"fmt"
	"testing"

	"github.com/jiushengTech/common/draw"
)

func TestPathDrawing(t *testing.T) {
	// 图片路径 - 使用在线URL
	imageURL := "http://suuqjbby1.hn-bkt.clouddn.com/tryon/origin/uploads/2025-04-22/20250420214131.jpg"

	// 创建一个"轮廓"路径
	outlinePath := draw.NewPath("outline")

	// 创建几个矩形，表示轮廓
	rect1 := draw.NewRectangle(
		draw.Point{X: 100, Y: 100},
		draw.Point{X: 200, Y: 200},
		draw.WithColor(draw.ColorBlue),
		draw.WithLineWidth(3.0),
	)

	rect2 := draw.NewRectangle(
		draw.Point{X: 300, Y: 150},
		draw.Point{X: 400, Y: 250},
		draw.WithColor(draw.ColorRed),
		draw.WithLineWidth(3.0),
	)

	// 添加矩形到轮廓路径
	outlinePath.AddShapes([]draw.Shape{rect1, rect2})

	// 创建一个"测量线"路径
	measurePath := draw.NewPath("measure")

	// 创建水平和垂直线，表示测量
	xpoints := []draw.Point{
		{X: 279, Y: 0},
		{X: 380, Y: 0},
		{X: 494, Y: 0},
		{X: 626, Y: 0},
	}
	xvalues := []float64{0.33, 0.45, 0.67}

	ypoints := []draw.Point{
		{X: 0, Y: 200},
		{X: 0, Y: 300},
		{X: 0, Y: 400},
	}
	yvalues := []float64{0.25, 0.50}

	vline := draw.NewVerticalLine(xpoints, xvalues, draw.WithColor(draw.ColorWhite))
	hline := draw.NewHorizontalLine(ypoints, yvalues, draw.WithColor(draw.ColorYellow))

	// 添加线到测量路径
	measurePath.AddShapes([]draw.Shape{vline, hline})

	// 创建一个"标记"路径
	markerPath := draw.NewPath("markers")

	// 创建一些圆形标记
	circle1 := draw.NewCircle(
		draw.Point{X: 500, Y: 300},
		20,
		draw.WithColor(draw.ColorGreen),
		draw.WithFill(true),
	)

	circle2 := draw.NewCircle(
		draw.Point{X: 650, Y: 300},
		20,
		draw.WithColor(draw.ColorRed),
		draw.WithFill(true),
	)

	// 添加圆形到标记路径
	markerPath.AddShapes([]draw.Shape{circle1, circle2})

	// 创建处理器并添加所有路径
	processor := draw.NewImageProcessor(
		imageURL,
		draw.WithTimeBasedName(),
		draw.WithOutputDir("path_results"),
		draw.WithShape(outlinePath),
		draw.WithShape(measurePath),
		draw.WithShape(markerPath),
	)

	// 处理图像
	if err := processor.Process(); err != nil {
		t.Errorf("绘制路径错误: %v", err)
	} else {
		fmt.Printf("成功生成包含多个路径的图片: %s/%s\n", processor.OutputDir, processor.Output)
	}

	// 测试单独开关某个路径
	// 隐藏测量路径
	measurePath.SetVisible(false)

	processor2 := draw.NewImageProcessor(
		imageURL,
		draw.WithTimeBasedName(),
		draw.WithOutputDir("path_results"),
		draw.WithShape(outlinePath),
		draw.WithShape(measurePath), // 这个路径被设置为不可见
		draw.WithShape(markerPath),
	)

	if err := processor2.Process(); err != nil {
		t.Errorf("绘制部分路径错误: %v", err)
	} else {
		fmt.Printf("成功生成隐藏测量路径的图片: %s/%s\n", processor2.OutputDir, processor2.Output)
	}
}
