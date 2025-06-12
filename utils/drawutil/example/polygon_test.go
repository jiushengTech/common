package example

import (
	"testing"

	"github.com/jiushengTech/common/utils/drawutil/colorx"
	"github.com/jiushengTech/common/utils/drawutil/processor"
	"github.com/jiushengTech/common/utils/drawutil/shape/base"
	"github.com/jiushengTech/common/utils/drawutil/shape/primitives/polygon"
)

// TestPolygonWithProcessor 测试使用ImageProcessor绘制多边形
func TestPolygonWithProcessor(t *testing.T) {
	// 图片URL
	imageURL := "http://139.224.59.6:9000/aesthetica/uploads/2025-06-10/149381445582850.jpg"

	// 创建第一个多边形（红色三角形，不填充）
	poly1 := polygon.New(false, // 不填充
		base.WithPoints([]*base.Point{
			{X: 200, Y: 200}, // 使用绝对坐标，ImageProcessor会根据图像尺寸进行调整
			{X: 400, Y: 300},
			{X: 100, Y: 400},
		}),
		base.WithColor(colorx.Red),
		base.WithLineWidth(3.0),
	)

	// 创建第二个多边形（蓝色三角形，填充）
	poly2 := polygon.New(true, // 填充
		base.WithPoints([]*base.Point{
			{X: 250, Y: 250}, // 使用绝对坐标，ImageProcessor会根据图像尺寸进行调整
			{X: 450, Y: 350},
			{X: 150, Y: 450},
		}),
		base.WithColor(colorx.Blue),
		base.WithLineWidth(3.0),
	)

	// 创建ImageProcessor，处理URL图片和添加图形
	imgProcessor := processor.NewImageProcessor(
		imageURL,
		processor.WithOutputDir("."), // 输出到example目录
		processor.WithOutputFormat(processor.FormatPNG), // 输出PNG格式
		processor.WithShape(poly1),                      // 添加第一个多边形
		processor.WithShape(poly2),                      // 添加第二个多边形
	)

	// 处理图像并保存
	outputPath, err := imgProcessor.Process()
	if err != nil {
		t.Fatalf("处理图像失败: %v", err)
	}

	t.Logf("图片已成功保存到: %s", outputPath)
}
