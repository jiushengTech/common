package main

import (
	"fmt"
	"os"
	"path/filepath"

	"linedraw"
)

func main() {
	// 获取当前执行目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前目录失败: %v\n", err)
		return
	}

	// 构建图片的绝对路径
	imagePath := filepath.Join(currentDir, "example.jpg")
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		fmt.Printf("找不到图片文件: %s\n", imagePath)
		return
	}

	// 示例1：只绘制竖线
	fmt.Println("示例1：绘制竖线...")

	// 创建竖线
	xpoints := []linedraw.Point{
		{X: 279, Y: 0},
		{X: 380, Y: 0},
		{X: 494, Y: 0},
		{X: 626, Y: 0},
		{X: 735, Y: 0},
		{X: 799, Y: 0},
	}
	xvalues := []float64{0.33, 0.45, 0.67, 0.82, 0.91}

	// 使用选项模式创建竖线
	vline := linedraw.NewVerticalLine(
		xpoints,
		xvalues,
		linedraw.WithColor(linedraw.ColorWhite), // 设置线条颜色为白色
		linedraw.WithLineWidth(2.5),             // 设置线宽
	)

	// 创建并处理图像
	processor1 := linedraw.NewImageProcessor(
		imagePath,
		linedraw.WithTimeBasedName(), // 使用时间作为文件名
		linedraw.WithLine(vline),     // 添加线条
	)

	if err := processor1.Process(); err != nil {
		fmt.Printf("绘制竖线错误: %v\n", err)
	} else {
		fmt.Printf("竖线绘制成功: %s/%s\n", processor1.OutputDir, processor1.Output)
	}

	// 示例2：只绘制横线
	fmt.Println("\n示例2：绘制横线...")

	// 创建横线
	ypoints := []linedraw.Point{
		{X: 0, Y: 200},
		{X: 0, Y: 300},
		{X: 0, Y: 400},
		{X: 0, Y: 500},
		{X: 0, Y: 600},
	}
	yvalues := []float64{0.25, 0.50, 0.75, 0.95}

	// 直接使用选项模式创建新的处理器
	processor2 := linedraw.NewImageProcessor(
		imagePath,
		linedraw.WithTimeBasedName(), // 使用时间作为文件名
		linedraw.WithLine(
			linedraw.NewHorizontalLine(
				ypoints,
				yvalues,
				linedraw.WithColor(linedraw.ColorBlue),
			),
		),
	)

	if err := processor2.Process(); err != nil {
		fmt.Printf("绘制横线错误: %v\n", err)
	} else {
		fmt.Printf("横线绘制成功: %s/%s\n", processor2.OutputDir, processor2.Output)
	}

	// 示例3：同时绘制竖线和横线
	fmt.Println("\n示例3：同时绘制竖线和横线...")

	// 创建并处理图像
	processor3 := linedraw.NewImageProcessor(
		imagePath,
		linedraw.WithTimeBasedName(), // 使用时间作为文件名
		linedraw.WithOutputDir("combined_results"), // 自定义输出目录
		linedraw.WithLines([]linedraw.Line{vline,
			linedraw.NewHorizontalLine(
				ypoints,
				yvalues,
				linedraw.WithColor(linedraw.ColorRed),
			),
		}),
	)

	if err := processor3.Process(); err != nil {
		fmt.Printf("绘制错误: %v\n", err)
	} else {
		fmt.Printf("成功生成图片: %s/%s\n", processor3.OutputDir, processor3.Output)
	}
}
