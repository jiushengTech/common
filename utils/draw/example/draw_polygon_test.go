package example

import (
	"fmt"
	"github.com/jiushengTech/common/utils/draw/shape/base"
	"testing"

	"github.com/jiushengTech/common/utils/draw"
)

func TestPolygonDrawing(t *testing.T) {
	// 图片路径 - 使用在线URL
	imageURL := "http://suuqjbby1.hn-bkt.clouddn.com/tryon/origin/uploads/2025-04-22/20250420214131.jpg"

	// 示例1：绘制一个普通的多边形
	fmt.Println("示例1：绘制普通多边形...")

	// 创建一个三角形
	trianglePoints := []*draw.Point{
		{X: 100, Y: 100},
		{X: 200, Y: 100},
		{X: 150, Y: 200},
	}

	triangle := draw.NewPolygon(
		trianglePoints,
		draw.WithColor(base.ColorRed),
		draw.WithLineWidth(3.0),
	)

	// 创建一个五边形
	pentagonPoints := []*draw.Point{
		{X: 300, Y: 150},
		{X: 350, Y: 100},
		{X: 400, Y: 150},
		{X: 380, Y: 200},
		{X: 320, Y: 200},
	}

	pentagon := draw.NewPolygon(
		pentagonPoints,
		draw.WithColor(base.ColorBlue),
		draw.WithPolygonFill(false), // 填充
	)

	// 创建一个不规则多边形
	irregularPoints := []*draw.Point{
		{X: 500, Y: 100},
		{X: 600, Y: 150},
		{X: 650, Y: 220},
		{X: 570, Y: 250},
		{X: 520, Y: 230},
		{X: 480, Y: 180},
	}

	irregular := draw.NewPolygon(
		irregularPoints,
		draw.WithColor(base.ColorGreen),
		draw.WithLineWidth(2.0),
	)

	// 处理图像
	processor := draw.NewImageProcessor(
		imageURL,
		draw.WithTimeBasedName(),
		draw.WithOutputDir("polygon_results"),
		draw.WithShape(triangle),
		draw.WithShape(pentagon),
		draw.WithShape(irregular),
	)

	absPath, err := processor.Process()
	if err != nil {
		t.Errorf("绘制多边形错误: %v", err)
	} else {
		fmt.Printf("成功生成多边形图片，绝对路径: %s\n", absPath)
	}

	// 示例2：创建一个复杂形状（星形）
	fmt.Println("示例2：绘制星形...")

	starPoints := []*draw.Point{
		{X: 400, Y: 300}, // 顶点
		{X: 425, Y: 350},
		{X: 475, Y: 350}, // 右上角
		{X: 435, Y: 375},
		{X: 450, Y: 425}, // 右下角
		{X: 400, Y: 400},
		{X: 350, Y: 425}, // 左下角
		{X: 365, Y: 375},
		{X: 325, Y: 350}, // 左上角
		{X: 375, Y: 350},
	}

	star := draw.NewPolygon(
		starPoints,
		draw.WithColor(base.ColorYellow),
		draw.WithPolygonFill(false),
	)

	// 处理图像
	processor2 := draw.NewImageProcessor(
		imageURL,
		draw.WithTimeBasedName(),
		draw.WithOutputDir("polygon_results"),
		draw.WithShape(star),
	)

	absPath2, err := processor2.Process()
	if err != nil {
		t.Errorf("绘制星形错误: %v", err)
	} else {
		fmt.Printf("成功生成星形图片，绝对路径: %s\n", absPath2)
	}
}
