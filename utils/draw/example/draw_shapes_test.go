package main

import (
	"fmt"
	"testing"

	"github.com/fogleman/gg"

	"github.com/jiushengTech/common/draw"
)

func TestShapesDrawing(t *testing.T) {
	// 图片路径 - 使用在线URL
	imageURL := "http://suuqjbby1.hn-bkt.clouddn.com/tryon/origin/uploads/2025-04-22/20250420214131.jpg"

	// 示例1：绘制线条
	fmt.Println("示例1：绘制线条...")

	// 创建竖线点和值
	xpoints := []draw.Point{
		{X: 279, Y: 0},
		{X: 380, Y: 0},
		{X: 494, Y: 0},
		{X: 626, Y: 0},
		{X: 735, Y: 0},
		{X: 799, Y: 0},
	}
	xvalues := []float64{0.33, 0.45, 0.67, 0.82, 0.91}

	// 创建垂直线
	vline := draw.NewVerticalLine(xpoints, xvalues, draw.WithColor(draw.ColorWhite))

	// 创建水平线
	ypoints := []draw.Point{
		{X: 0, Y: 200},
		{X: 0, Y: 300},
		{X: 0, Y: 400},
		{X: 0, Y: 500},
	}
	yvalues := []float64{0.25, 0.50, 0.75}

	hline := draw.NewHorizontalLine(ypoints, yvalues, draw.WithColor(draw.ColorYellow))

	// 示例2：绘制矩形
	fmt.Println("示例2：绘制矩形...")

	// 创建矩形
	rect1 := draw.NewRectangle(
		draw.Point{X: 100, Y: 100},
		draw.Point{X: 200, Y: 200},
		draw.WithColor(draw.ColorBlue),
		draw.WithLineWidth(3.0),
	)

	// 创建填充矩形
	rect2 := draw.NewRectangle(
		draw.Point{X: 300, Y: 150},
		draw.Point{X: 400, Y: 250},
		draw.WithColor(draw.ColorRed),
		draw.WithFill(true),
	)

	// 示例3：绘制圆形
	fmt.Println("示例3：绘制圆形...")

	// 创建圆形
	circle1 := draw.NewCircle(
		draw.Point{X: 500, Y: 300},
		50,
		draw.WithColor(draw.ColorGreen),
		draw.WithLineWidth(2.5),
	)

	// 创建填充圆形
	circle2 := draw.NewCircle(
		draw.Point{X: 650, Y: 300},
		30,
		draw.WithColor(draw.ColorBlue),
		draw.WithFill(true),
	)

	// 将所有图形放入切片
	shapes := []draw.Shape{vline, hline, rect1, rect2, circle1, circle2}

	// 创建并处理图像
	processor := draw.NewImageProcessor(
		imageURL,
		draw.WithTimeBasedName(),
		draw.WithOutputDir("combined_results"),
		draw.WithShapes(shapes),
	)

	if err := processor.Process(); err != nil {
		t.Errorf("绘制图形错误: %v", err)
	} else {
		fmt.Printf("成功生成包含多种图形的图片: %s/%s\n", processor.OutputDir, processor.Output)
	}

	// 示例4：使用工厂模式创建图形
	fmt.Println("示例4：使用工厂模式创建图形...")

	// 通过类型名称创建图形
	circle, ok := draw.NewShape("circle",
		draw.WithPoints([]draw.Point{{X: 400, Y: 400}}),
		draw.WithRadius(60),
		draw.WithColor(draw.ColorRed),
		draw.WithFill(true),
	)

	if !ok {
		t.Error("创建圆形失败")
	}

	rectangle, ok := draw.NewShape("rectangle",
		draw.WithPoints([]draw.Point{{X: 200, Y: 350}, {X: 350, Y: 450}}),
		draw.WithColor(draw.ColorGreen),
		draw.WithFill(true),
	)

	if !ok {
		t.Error("创建矩形失败")
	}

	processor2 := draw.NewImageProcessor(
		imageURL,
		draw.WithTimeBasedName(),
		draw.WithOutputDir("factory_results"),
		draw.WithShape(circle),
		draw.WithShape(rectangle),
	)

	if err := processor2.Process(); err != nil {
		t.Errorf("绘制工厂创建的图形错误: %v", err)
	} else {
		fmt.Printf("成功生成工厂创建的图形: %s/%s\n", processor2.OutputDir, processor2.Output)
	}

	// 示例5：使用新增的格式和质量选项
	fmt.Println("示例5：使用JPEG格式和自定义质量...")

	// 创建一个处理器，输出JPEG格式的图像，设置75%的质量
	processor3 := draw.NewImageProcessor(
		imageURL,
		draw.WithOutputDir("jpeg_results"),
		draw.WithOutputFormat(draw.FormatJPEG),
		draw.WithJpegQuality(75),
		draw.WithTimeBasedName(),
		draw.WithShape(circle1),
		draw.WithShape(rect1),
	)

	if err := processor3.Process(); err != nil {
		t.Errorf("生成JPEG图像失败: %v", err)
	} else {
		fmt.Printf("成功生成JPEG格式图像: %s/%s\n", processor3.OutputDir, processor3.Output)
	}

	// 示例6：使用预处理和后处理函数
	fmt.Println("示例6：使用预处理和后处理函数...")

	// 创建一个处理器，添加预处理和后处理函数
	processor4 := draw.NewImageProcessor(
		imageURL,
		draw.WithOutputDir("processed_results"),
		draw.WithTimeBasedName(),
		// 预处理函数 - 在绘制形状前应用灰度效果
		draw.WithPreProcess(func(dc *gg.Context, width, height float64) error {
			// 遍历每个像素并转换为灰度
			img := dc.Image()
			for y := 0; y < int(height); y++ {
				for x := 0; x < int(width); x++ {
					c := img.At(x, y)
					r, g, b, _ := c.RGBA()
					// 简单的灰度公式
					gray := uint8((0.3*float64(r) + 0.59*float64(g) + 0.11*float64(b)) / 256.0)
					dc.SetRGB(float64(gray)/255.0, float64(gray)/255.0, float64(gray)/255.0)
					dc.SetPixel(x, y)
				}
			}
			return nil
		}),
		draw.WithShape(circle2),
		draw.WithShape(rect2),
		// 后处理函数 - 添加文本水印
		draw.WithPostProcess(func(dc *gg.Context, width, height float64) error {
			// 水印文本
			text := "© jiushengTech"

			// 尝试直接设置字体大小而不是加载特定字体
			dc.SetLineWidth(1)
			// 不使用LoadFontFace，而是直接绘制

			// 计算水印位置（右下角，但留出一些边距）
			margin := 20.0
			x := width - margin
			y := height - margin

			// 绘制一个半透明背景以确保水印可见
			dc.SetRGBA(0, 0, 0, 0.5) // 半透明黑色背景
			textWidth := 180.0       // 估计文本宽度
			textHeight := 30.0       // 估计文本高度
			dc.DrawRoundedRectangle(x-textWidth, y-textHeight, textWidth, textHeight, 5)
			dc.Fill()

			// 设置文本前景色为白色
			dc.SetRGB(1, 1, 1)

			// 绘制多行简单文本，不依赖字体
			simpleText := func(text string, x, y float64) {
				// 绘制文本的简单方法
				dc.DrawString(text, x-textWidth+10, y-15)
			}

			simpleText(text, x, y)

			return nil
		}),
	)

	if err := processor4.Process(); err != nil {
		t.Errorf("使用处理函数生成图像失败: %v", err)
	} else {
		fmt.Printf("成功生成处理后的图像: %s/%s\n", processor4.OutputDir, processor4.Output)
	}

	// 清理临时文件
	draw.CleanupAllTempFiles()
}
