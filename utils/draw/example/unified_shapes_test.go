package example

import (
	"testing"

	"github.com/jiushengTech/common/utils/draw/colorx"
	"github.com/jiushengTech/common/utils/draw/processor"
	"github.com/jiushengTech/common/utils/draw/shape/base"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/circle"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/line"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/polygon"
)

// TestUnifiedShapeInterface 测试统一的图形接口
func TestUnifiedShapeInterface(t *testing.T) {
	// 图片URL
	imageURL := "http://139.224.59.6:9000/aesthetica/uploads/2025-06-10/149381445582850.jpg"

	// 1. 创建多边形（红色三角形，不填充）
	poly := polygon.New(false, // 不填充
		base.WithPoints([]*base.Point{
			{X: 100, Y: 100},
			{X: 300, Y: 150},
			{X: 200, Y: 300},
		}),
		base.WithColor(colorx.Red),
		base.WithLineWidth(3.0),
	)

	// 2. 创建矩形（用多边形实现，蓝色，填充）
	rect := polygon.New(true, // 填充
		base.WithPoints([]*base.Point{
			{X: 400, Y: 100}, // 左上角
			{X: 600, Y: 100}, // 右上角
			{X: 600, Y: 250}, // 右下角
			{X: 400, Y: 250}, // 左下角
		}),
		base.WithColor(colorx.Blue),
		base.WithLineWidth(2.0),
	)

	// 3. 创建圆形（绿色，不填充）
	circleShape := circle.New(50.0, false, // 半径50，不填充
		base.WithPoints([]*base.Point{
			{X: 150, Y: 400}, // 圆心
		}),
		base.WithColor(colorx.Green),
		base.WithLineWidth(4.0),
	)

	// 4. 创建垂直线条（黄色，带数值）
	vertLine := line.New(line.Vertical, 0.8, // 垂直线，文字位置在80%高度
		base.WithPoints([]*base.Point{
			{X: 350, Y: 0},
			{X: 450, Y: 0},
			{X: 550, Y: 0},
		}),
		base.WithColor(colorx.Yellow),
		base.WithLineWidth(3.0),
	)
	// 设置线条值
	vertLine.SetValues([]float64{15.5, 25.8})

	// 5. 创建水平线条（橙色，带数值）
	horizLine := line.New(line.Horizontal, 0.2, // 水平线，文字位置在20%宽度
		base.WithPoints([]*base.Point{
			{X: 0, Y: 450},
			{X: 0, Y: 550},
		}),
		base.WithColor(colorx.Orange),
		base.WithLineWidth(2.5),
	)
	// 设置线条值
	horizLine.SetValues([]float64{12.3})

	// 创建图像处理器并添加所有图形
	imgProcessor := processor.NewImageProcessor(
		imageURL,
		processor.WithOutputDir("."),
		processor.WithOutputFormat(processor.FormatPNG),
		processor.WithShape(poly),
		processor.WithShape(rect),
		processor.WithShape(circleShape),
		processor.WithShape(vertLine),
		processor.WithShape(horizLine),
	)

	// 处理图像并保存
	outputPath, err := imgProcessor.Process()
	if err != nil {
		t.Fatalf("处理图像失败: %v", err)
	}

	t.Logf("统一接口图形演示已保存到: %s", outputPath)
}

// TestUnifiedInterfaceBenefits 测试统一接口的优势
func TestUnifiedInterfaceBenefits(t *testing.T) {
	t.Log("=== 统一接口的优势 ===")

	t.Log("1. 一致的构造函数模式：")
	t.Log("   - polygon.New(fill, options...) // 支持三角形、矩形等所有多边形")
	t.Log("   - circle.New(radius, fill, options...)")
	t.Log("   - line.New(type, textPosition, options...)")

	t.Log("2. 统一的选项系统：")
	t.Log("   - 所有图形都使用 base.Option")
	t.Log("   - base.WithPoints() 设置坐标点")
	t.Log("   - base.WithColor() 设置颜色")
	t.Log("   - base.WithLineWidth() 设置线宽")

	t.Log("3. 简化的代码结构：")
	t.Log("   - 删除了所有重复的选项适配器")
	t.Log("   - 删除了复杂的多重构造函数")
	t.Log("   - 统一的参数传递方式")
	t.Log("   - 矩形作为多边形的特例，减少了冗余代码")

	t.Log("4. 易于维护和扩展：")
	t.Log("   - 新增通用选项只需在 base 包中添加")
	t.Log("   - 图形特有属性作为构造函数参数")
	t.Log("   - 清晰的职责分离")
}
