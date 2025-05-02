package main

import (
	"fmt"
	"testing"

	"github.com/jiushengTech/common/utils/draw"
)

func TestDrawLines(t *testing.T) {
	// 测试图片路径 - 使用在线URL
	imageURL := "http://suuqjbby1.hn-bkt.clouddn.com/tryon/origin/uploads/2025-04-22/20250420214131.jpg"

	// 测试用例1：绘制垂直线 - 文本在中间位置
	t.Run("绘制垂直线-中间文本", func(t *testing.T) {
		// 创建竖线点和值
		xpoints := []draw.Point{
			{X: 100, Y: 0}, // 第一条竖线
			{X: 200, Y: 0}, // 第二条竖线
			{X: 300, Y: 0}, // 第三条竖线
			{X: 400, Y: 0}, // 第四条竖线
			{X: 500, Y: 0}, // 第五条竖线
		}
		// 这些值将显示在相邻两条线之间
		xvalues := []float64{0.25, 0.50, 0.75, 1.0}

		// 创建垂直线 - 设置文本在正中间 (0.5)
		vline := draw.NewVerticalLine(xpoints, xvalues,
			draw.WithColor(draw.ColorRed),
			draw.WithLineWidth(2.0),
			draw.WithTextPosition(0.5), // 文本位置在正中间
		)

		// 创建图像处理器并添加图形
		processor := draw.NewImageProcessor(
			imageURL,
			draw.WithOutputDir("lines_test"),
			draw.WithOutputName("vertical_lines_center.png"),
			draw.WithShape(vline),
		)
		// 处理图像
		absPath, err := processor.Process()
		if err != nil {
			t.Errorf("绘制垂直线错误: %v", err)
		} else {
			fmt.Printf("成功生成垂直线图片(中间文本)，绝对路径: %s\n", absPath)
		}
	})

	// 测试用例2：绘制垂直线 - 文本在顶部
	t.Run("绘制垂直线-顶部文本", func(t *testing.T) {
		// 创建竖线点和值
		xpoints := []draw.Point{
			{X: 100, Y: 0}, // 第一条竖线
			{X: 200, Y: 0}, // 第二条竖线
			{X: 300, Y: 0}, // 第三条竖线
			{X: 400, Y: 0}, // 第四条竖线
			{X: 500, Y: 0}, // 第五条竖线
		}
		// 这些值将显示在相邻两条线之间
		xvalues := []float64{0.25, 0.50, 0.75, 1.0}

		// 创建垂直线 - 设置文本在顶部 (0.1)
		vline := draw.NewVerticalLine(xpoints, xvalues,
			draw.WithColor(draw.ColorRed),
			draw.WithLineWidth(2.0),
			draw.WithTextPosition(0.1), // 文本位置在顶部
		)

		// 创建图像处理器并添加图形
		processor := draw.NewImageProcessor(
			imageURL,
			draw.WithOutputDir("lines_test"),
			draw.WithOutputName("vertical_lines_top.png"),
			draw.WithShape(vline),
		)

		// 处理图像
		absPath, err := processor.Process()
		if err != nil {
			t.Errorf("绘制垂直线错误: %v", err)
		} else {
			fmt.Printf("成功生成垂直线图片(顶部文本)，绝对路径: %s\n", absPath)
		}
	})

	// 测试用例3：绘制水平线 - 文本在中间
	t.Run("绘制水平线-中间文本", func(t *testing.T) {
		// 创建水平线点和值
		ypoints := []draw.Point{
			{X: 0, Y: 100}, // 第一条水平线
			{X: 0, Y: 200}, // 第二条水平线
			{X: 0, Y: 300}, // 第三条水平线
			{X: 0, Y: 400}, // 第四条水平线
			{X: 0, Y: 500}, // 第五条水平线
		}
		// 这些值将显示在相邻两条线之间
		yvalues := []float64{0.25, 0.50, 0.75, 1.0}

		// 创建水平线 - 设置文本在正中间
		hline := draw.NewHorizontalLine(ypoints, yvalues,
			draw.WithColor(draw.ColorBlue),
			draw.WithLineWidth(2.0),
			draw.WithTextPosition(0.5), // 文本位置在正中间
		)

		// 创建图像处理器并添加图形
		processor := draw.NewImageProcessor(
			imageURL,
			draw.WithOutputDir("lines_test"),
			draw.WithOutputName("horizontal_lines_center.png"),
			draw.WithShape(hline),
		)

		// 处理图像
		absPath, err := processor.Process()
		if err != nil {
			t.Errorf("绘制水平线错误: %v", err)
		} else {
			fmt.Printf("成功生成水平线图片(中间文本)，绝对路径: %s\n", absPath)
		}
	})

	// 测试用例4：绘制水平线 - 文本在右侧
	t.Run("绘制水平线-右侧文本", func(t *testing.T) {
		// 创建水平线点和值
		ypoints := []draw.Point{
			{X: 0, Y: 100}, // 第一条水平线
			{X: 0, Y: 200}, // 第二条水平线
			{X: 0, Y: 300}, // 第三条水平线
			{X: 0, Y: 400}, // 第四条水平线
			{X: 0, Y: 500}, // 第五条水平线
		}
		// 这些值将显示在相邻两条线之间
		yvalues := []float64{0.25, 0.50, 0.75, 1.0}

		// 创建水平线 - 设置文本在右侧
		hline := draw.NewHorizontalLine(ypoints, yvalues,
			draw.WithColor(draw.ColorBlue),
			draw.WithLineWidth(2.0),
			draw.WithTextPosition(0.9), // 文本位置在右侧
		)

		// 创建图像处理器并添加图形
		processor := draw.NewImageProcessor(
			imageURL,
			draw.WithOutputDir("lines_test"),
			draw.WithOutputName("horizontal_lines_right.png"),
			draw.WithShape(hline),
		)

		// 处理图像
		absPath, err := processor.Process()
		if err != nil {
			t.Errorf("绘制水平线错误: %v", err)
		} else {
			fmt.Printf("成功生成水平线图片(右侧文本)，绝对路径: %s\n", absPath)
		}
	})

	// 测试用例5：同时绘制不同位置文本的水平线和垂直线
	t.Run("绘制不同文本位置的线条组合", func(t *testing.T) {
		// 垂直线 - 中间文本
		vline1 := draw.NewVerticalLine(
			[]draw.Point{{X: 150, Y: 0}, {X: 300, Y: 0}, {X: 450, Y: 0}},
			[]float64{0.25, 0.5},
			draw.WithColor(draw.ColorRed),
			draw.WithLineWidth(2.0),
			draw.WithTextPosition(0.5), // 中间文本
		)

		// 垂直线 - 底部文本
		vline2 := draw.NewVerticalLine(
			[]draw.Point{{X: 600, Y: 0}, {X: 750, Y: 0}},
			[]float64{0.75},
			draw.WithColor(draw.ColorGreen),
			draw.WithLineWidth(2.0),
			draw.WithTextPosition(0.85), // 底部文本
		)

		// 水平线 - 左侧文本
		hline1 := draw.NewHorizontalLine(
			[]draw.Point{{X: 0, Y: 150}, {X: 0, Y: 300}, {X: 0, Y: 450}},
			[]float64{0.25, 0.5},
			draw.WithColor(draw.ColorBlue),
			draw.WithLineWidth(2.0),
			draw.WithTextPosition(0.1), // 左侧文本
		)

		// 水平线 - 右侧文本
		hline2 := draw.NewHorizontalLine(
			[]draw.Point{{X: 0, Y: 600}, {X: 0, Y: 750}},
			[]float64{0.75},
			draw.WithColor(draw.ColorYellow),
			draw.WithLineWidth(2.0),
			draw.WithTextPosition(0.9), // 右侧文本
		)

		// 创建图像处理器并添加所有图形
		processor := draw.NewImageProcessor(
			imageURL,
			draw.WithOutputDir("lines_test"),
			draw.WithOutputName("combined_text_positions.png"),
			draw.WithShape(vline1),
			draw.WithShape(vline2),
			draw.WithShape(hline1),
			draw.WithShape(hline2),
		)

		// 处理图像
		absPath, err := processor.Process()
		if err != nil {
			t.Errorf("绘制组合线条错误: %v", err)
		} else {
			fmt.Printf("成功生成组合线条图片，绝对路径: %s\n", absPath)
		}
	})
}
