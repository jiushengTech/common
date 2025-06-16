package example

import (
	"testing"

	"github.com/jiushengTech/common/utils/drawutil/colorx"
	"github.com/jiushengTech/common/utils/drawutil/processor"
	"github.com/jiushengTech/common/utils/drawutil/shape/base"
	"github.com/jiushengTech/common/utils/drawutil/shape/primitives/line"
)

// TestLineWithNewOptionPattern 测试新的选项模式
func TestLineWithNewOptionPattern(t *testing.T) {
	// 图片URL
	imageURL := "http://139.224.59.6:9000/aesthetica/uploads/2025-06-10/149381445582850.jpg"

	// 1. 创建垂直线条（使用默认值）
	vertLine1 := line.New(
		line.WithPoints([]*base.Point{
			{X: 100, Y: 0},
			{X: 200, Y: 0},
		}),
		line.WithColor(colorx.Red),
	)

	// 2. 创建垂直线条（完整配置）
	vertLine2 := line.New(
		line.WithType(line.Vertical),
		line.WithTextPosition(0.3),
		line.WithValues([]float64{10.5, 20.8, 15.2}),
		line.WithPoints([]*base.Point{
			{X: 300, Y: 0},
			{X: 400, Y: 0},
			{X: 500, Y: 0},
			{X: 600, Y: 0},
		}),
		line.WithColor(colorx.Blue),
		line.WithLineWidth(4.0),
	)

	// 3. 创建水平线条
	horizLine := line.New(
		line.WithType(line.Horizontal),
		line.WithTextPosition(0.7),
		line.WithValues([]float64{25.3, 18.7}),
		line.WithPoints([]*base.Point{
			{X: 0, Y: 150},
			{X: 0, Y: 250},
			{X: 0, Y: 350},
		}),
		line.WithColor(colorx.Green),
		line.WithLineWidth(3.0),
	)

	// 4. 创建简单水平线（最小配置）
	simpleHorizLine := line.New(
		line.WithType(line.Horizontal),
		line.WithPoints([]*base.Point{
			{X: 0, Y: 450},
			{X: 0, Y: 550},
		}),
		line.WithColor(colorx.Orange),
	)

	// 创建图像处理器并添加所有图形
	imgProcessor := processor.NewImageProcessor(
		imageURL,
		processor.WithOutputDir("."),
		processor.WithOutputFormat(processor.FormatPNG),
		processor.WithShape(vertLine1),
		processor.WithShape(vertLine2),
		processor.WithShape(horizLine),
		processor.WithShape(simpleHorizLine),
	)

	// 处理图像并保存
	outputPath, err := imgProcessor.Process()
	if err != nil {
		t.Fatalf("处理图像失败: %v", err)
	}

	t.Logf("线条选项模式演示已保存到: %s", outputPath)
}

// TestLineOptionsBenefits 测试线条选项模式的优势
func TestLineOptionsBenefits(t *testing.T) {
	t.Log("=== 优雅的线条选项模式优势 ===")

	t.Log("1. 统一的API接口：")
	t.Log("   - line.New(options...) // 单一可变参数")
	t.Log("   - 所有选项都使用相同的LineOption类型")

	t.Log("2. 直观的选项组合：")
	t.Log("   - line.WithType() 设置线条类型")
	t.Log("   - line.WithValues() 设置值集合")
	t.Log("   - line.WithTextPosition() 设置文本位置")
	t.Log("   - line.WithPoints() 设置坐标点")
	t.Log("   - line.WithColor() 设置颜色")
	t.Log("   - line.WithLineWidth() 设置线宽")

	t.Log("3. 简洁的使用方式：")
	t.Log("   - 所有配置在一个函数调用中完成")
	t.Log("   - 选项顺序任意，灵活组合")
	t.Log("   - 不需要分离线条选项和基础选项")

	t.Log("4. 类型安全和一致性：")
	t.Log("   - 编译时检查所有选项类型")
	t.Log("   - 统一的选项函数签名")
	t.Log("   - 避免了参数传递错误")

	t.Log("5. 易于维护和扩展：")
	t.Log("   - 新增选项只需添加新的WithXxx函数")
	t.Log("   - 所有选项都在options.go中集中管理")
	t.Log("   - 向后兼容，不破坏现有代码")
}
