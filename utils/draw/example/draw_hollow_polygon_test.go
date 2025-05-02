package main

import (
	"fmt"
	"testing"

	"github.com/jiushengTech/common/utils/draw"
)

func TestHollowPolygonDrawing(t *testing.T) {
	// 图片路径 - 使用在线URL
	imageURL := "http://suuqjbby1.hn-bkt.clouddn.com/tryon/origin/uploads/2025-04-22/20250420214131.jpg"

	// 示例1：创建一个简单的镂空多边形（大方形中镂空一个小方形）
	fmt.Println("示例1：创建镂空矩形...")

	// 外部多边形（大矩形）
	outerPoints := []draw.Point{
		{X: 100, Y: 100}, // 左上
		{X: 400, Y: 100}, // 右上
		{X: 400, Y: 400}, // 右下
		{X: 100, Y: 400}, // 左下
	}

	// 内部多边形（小矩形）
	innerPoints := []draw.Point{
		{X: 200, Y: 200}, // 左上
		{X: 300, Y: 200}, // 右上
		{X: 300, Y: 300}, // 右下
		{X: 200, Y: 300}, // 左下
	}

	// 创建镂空多边形
	hollowRect := draw.NewHollowPolygon(
		outerPoints,
		innerPoints,
		draw.WithColor(draw.ColorGray),
		draw.WithHollowPolygonOpacity(0.5), // 不透明度为0.7
		draw.WithLineWidth(2.0),
	)

	// 创建处理器并处理图像
	processor := draw.NewImageProcessor(
		imageURL,
		draw.WithTimeBasedName(),
		draw.WithOutputDir("hollow_results"),
		draw.WithShape(hollowRect),
	)

	absPath, err := processor.Process()
	if err != nil {
		t.Errorf("绘制镂空矩形错误: %v", err)
	} else {
		fmt.Printf("成功生成镂空矩形图片，绝对路径: %s\n", absPath)
	}

	// 示例2：创建一个复杂的镂空多边形（不规则外形中镂空一个五边形）
	fmt.Println("示例2：创建复杂镂空多边形...")

	// 外部多边形（不规则形状）
	outerComplex := []draw.Point{
		{X: 150, Y: 150},
		{X: 350, Y: 120},
		{X: 450, Y: 250},
		{X: 400, Y: 380},
		{X: 250, Y: 450},
		{X: 100, Y: 350},
	}

	// 内部多边形（五边形）
	innerPentagon := []draw.Point{
		{X: 200, Y: 180},
		{X: 300, Y: 200},
		{X: 330, Y: 280},
		{X: 250, Y: 340},
		{X: 180, Y: 260},
	}

	// 创建镂空多边形
	hollowComplex := draw.NewHollowPolygon(
		outerComplex,
		innerPentagon,
		draw.WithColor(draw.ColorGreen),
		draw.WithHollowPolygonOpacity(0.5), // 不透明度为0.5
		draw.WithLineWidth(2.5),
	)

	// 创建处理器并处理图像
	processor2 := draw.NewImageProcessor(
		imageURL,
		draw.WithTimeBasedName(),
		draw.WithOutputDir("hollow_results"),
		draw.WithShape(hollowComplex),
	)

	absPath2, err := processor2.Process()
	if err != nil {
		t.Errorf("绘制复杂镂空多边形错误: %v", err)
	} else {
		fmt.Printf("成功生成复杂镂空多边形图片，绝对路径: %s\n", absPath2)
	}

	// 示例3：显示内部和外部多边形的边界
	fmt.Println("示例3：同时显示镂空多边形和普通多边形...")

	// 创建镂空多边形
	hollowWithBorder := draw.NewHollowPolygon(
		outerComplex,
		innerPentagon,
		draw.WithColor(draw.ColorRed),
		draw.WithHollowPolygonOpacity(0.3), // 低不透明度
		draw.WithLineWidth(3.0),
	)

	// 创建内部五边形的实体轮廓
	innerOutline := draw.NewPolygon(
		innerPentagon,
		draw.WithColor(draw.ColorYellow),
		draw.WithLineWidth(2.0),
	)

	// 创建处理器并处理图像
	processor3 := draw.NewImageProcessor(
		imageURL,
		draw.WithTimeBasedName(),
		draw.WithOutputDir("hollow_results"),
		draw.WithShape(hollowWithBorder), // 先绘制镂空区域
		draw.WithShape(innerOutline),     // 再绘制内部轮廓
	)

	absPath3, err := processor3.Process()
	if err != nil {
		t.Errorf("绘制镂空多边形和边界错误: %v", err)
	} else {
		fmt.Printf("成功生成带边界的镂空多边形图片，绝对路径: %s\n", absPath3)
	}
}
