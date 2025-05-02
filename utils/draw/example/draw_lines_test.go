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
		// 创建普通竖线
		xpoints := []draw.Point{
			{X: 279, Y: 0},
			{X: 380, Y: 0},
			{X: 494, Y: 0},
			{X: 626, Y: 0},
			{X: 735, Y: 0},
			{X: 799, Y: 0},
		}
		xvalues := []float64{0.33, 0.45, 0.67, 0.82, 0.91}

		vline := draw.NewLine(draw.VerticalLine, xpoints, xvalues,
			draw.WithColor(draw.ColorWhite))

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
		// 创建普通竖线
		xpoints := []draw.Point{
			{X: 279, Y: 0},
			{X: 380, Y: 0},
			{X: 494, Y: 0},
			{X: 626, Y: 0},
			{X: 735, Y: 0},
			{X: 799, Y: 0},
		}
		xvalues := []float64{0.33, 0.45, 0.67, 0.82, 0.91}

		vline := draw.NewLine(draw.VerticalLine, xpoints, xvalues,
			draw.WithColor(draw.ColorWhite))

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
		// 绘制简单横线 - 带颜色、线宽设置
		ypoints := []draw.Point{
			{X: 0, Y: 200},
			{X: 0, Y: 300},
			{X: 0, Y: 400},
			{X: 0, Y: 500},
		}
		yvalues := []float64{0.25, 0.50, 0.75}

		hline := draw.NewLine(draw.HorizontalLine, ypoints, yvalues,
			draw.WithColor(draw.ColorBlue),
			draw.WithLineWidth(2.0))

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
		// 绘制横线 - 带文本位置设置
		ypoints := []draw.Point{
			{X: 0, Y: 200},
			{X: 0, Y: 300},
			{X: 0, Y: 400},
			{X: 0, Y: 500},
		}
		yvalues := []float64{0.25, 0.50, 0.75}

		hline := draw.NewLine(draw.HorizontalLine, ypoints, yvalues,
			draw.WithColor(draw.ColorYellow),
			draw.WithTextPosition(0.1))

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
		// 竖线1 - 左侧
		vline1 := draw.NewLine(draw.VerticalLine,
			[]draw.Point{
				{X: 180, Y: 0},
				{X: 220, Y: 0},
				{X: 260, Y: 0},
			},
			[]float64{0.3, 0.7},
			draw.WithColor(draw.ColorRed),
			draw.WithTextPosition(0.1), // 偏左侧显示文本
		)

		// 竖线2 - 右侧
		vline2 := draw.NewLine(draw.VerticalLine,
			[]draw.Point{
				{X: 540, Y: 0},
				{X: 580, Y: 0},
				{X: 620, Y: 0},
			},
			[]float64{0.3, 0.7},
			draw.WithColor(draw.ColorBlue),
			draw.WithTextPosition(0.9), // 偏右侧显示文本
		)

		// 横线1 - 上方
		hline1 := draw.NewLine(draw.HorizontalLine,
			[]draw.Point{
				{X: 0, Y: 180},
				{X: 0, Y: 220},
				{X: 0, Y: 260},
			},
			[]float64{0.3, 0.7},
			draw.WithColor(draw.ColorGreen),
			draw.WithTextPosition(0.1), // 偏上方显示文本
		)

		// 横线2 - 下方
		hline2 := draw.NewLine(draw.HorizontalLine,
			[]draw.Point{
				{X: 0, Y: 540},
				{X: 0, Y: 580},
				{X: 0, Y: 620},
			},
			[]float64{0.3, 0.7},
			draw.WithColor(draw.ColorOrange),
			draw.WithTextPosition(0.9), // 偏下方显示文本
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
