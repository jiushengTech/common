package example

import (
	"fmt"
	"github.com/jiushengTech/common/utils/draw/shape/base"
	"testing"

	"github.com/jiushengTech/common/utils/draw"
)

func TestShapeGroupDrawing(t *testing.T) {
	// 图片路径 - 使用在线URL
	imageURL := "http://suuqjbby1.hn-bkt.clouddn.com/tryon/origin/uploads/2025-04-22/20250420214131.jpg"

	// 创建一个"轮廓"图形组
	outlineGroup := draw.NewShapeGroup("outline")

	// 创建几个矩形，表示轮廓
	rect1 := draw.NewRectangle(
		&draw.Point{X: 100, Y: 100},
		&draw.Point{X: 200, Y: 200},
		draw.WithColor(base.ColorBlue),
		draw.WithLineWidth(3.0),
	)

	rect2 := draw.NewRectangle(
		&draw.Point{X: 300, Y: 150},
		&draw.Point{X: 400, Y: 250},
		draw.WithColor(base.ColorRed),
		draw.WithLineWidth(3.0),
	)

	// 添加矩形到轮廓图形组
	outlineGroup.AddShapes([]draw.Shape{rect1, rect2})

	// 创建一个"测量线"图形组
	measureGroup := draw.NewShapeGroup("measure")

	// 创建水平和垂直线，表示测量
	xpoints := []*draw.Point{
		{X: 279, Y: 0},
		{X: 380, Y: 0},
		{X: 494, Y: 0},
		{X: 626, Y: 0},
	}
	xvalues := []float64{0.33, 0.45, 0.67}

	ypoints := []*draw.Point{
		{X: 0, Y: 200},
		{X: 0, Y: 300},
		{X: 0, Y: 400},
	}
	yvalues := []float64{0.25, 0.50}

	// 创建竖线组
	vline := draw.NewLine(draw.VerticalLine, xpoints, xvalues, draw.WithColor(base.ColorWhite))
	hline := draw.NewLine(draw.HorizontalLine, ypoints, yvalues, draw.WithColor(base.ColorYellow))

	// 添加线到测量图形组
	measureGroup.AddShapes([]draw.Shape{vline, hline})

	// 创建一个"标记"图形组
	markerGroup := draw.NewShapeGroup("markers")

	// 创建一些圆形标记
	circle1 := draw.NewCircle(
		&draw.Point{X: 500, Y: 300},
		20,
		draw.WithColor(base.ColorGreen),
		draw.WithFill(true),
	)

	circle2 := draw.NewCircle(
		&draw.Point{X: 650, Y: 300},
		20,
		draw.WithColor(base.ColorRed),
		draw.WithFill(true),
	)

	// 添加圆形到标记图形组
	markerGroup.AddShapes([]draw.Shape{circle1, circle2})

	// 创建处理器并添加所有图形组
	processor := draw.NewImageProcessor(
		imageURL,
		draw.WithTimeBasedName(),
		draw.WithOutputDir("group_results"),
		draw.WithShape(outlineGroup),
		draw.WithShape(measureGroup),
		draw.WithShape(markerGroup),
	)

	// 处理图像
	absPath, err := processor.Process()
	if err != nil {
		t.Errorf("绘制图形组错误: %v", err)
	} else {
		fmt.Printf("成功生成包含多个图形组的图片，绝对路径: %s\n", absPath)
	}

	// 测试单独开关某个图形组
	// 隐藏测量图形组
	measureGroup.SetVisible(false)

	processor2 := draw.NewImageProcessor(
		imageURL,
		draw.WithTimeBasedName(),
		draw.WithOutputDir("group_results"),
		draw.WithShape(outlineGroup),
		draw.WithShape(measureGroup), // 这个图形组被设置为不可见
		draw.WithShape(markerGroup),
	)

	absPath2, err := processor2.Process()
	if err != nil {
		t.Errorf("绘制部分图形组错误: %v", err)
	} else {
		fmt.Printf("成功生成隐藏测量图形组的图片，绝对路径: %s\n", absPath2)
	}
}
