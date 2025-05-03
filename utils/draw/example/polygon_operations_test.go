package example

import (
	"fmt"
	"math"
	"testing"

	"github.com/fogleman/gg"
	"github.com/jiushengTech/common/utils/draw"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/polygonops"
)

// 示例图片URL
const exampleImageURL = "http://suuqjbby1.hn-bkt.clouddn.com/tryon/origin/uploads/2025-04-22/20250420214131.jpg"

// TestPolygonOperationsBasic 测试基本多边形布尔运算（交集、差集）
func TestPolygonOperationsBasic(t *testing.T) {
	fmt.Println("测试基本多边形布尔运算（交集、差集）...")

	// 创建外部多边形A的点（矩形）
	polygonA := []*draw.Point{
		{X: 100, Y: 100}, // 左上
		{X: 400, Y: 100}, // 右上
		{X: 400, Y: 400}, // 右下
		{X: 100, Y: 400}, // 左下
	}

	// 创建内部多边形B的点（圆形 - 用多边形近似）
	polygonB := []*draw.Point{}
	centerX, centerY := 300.0, 250.0
	radius := 120.0

	// 创建36个点的圆形
	for i := 0; i < 36; i++ {
		angle := float64(i) * 2 * math.Pi / 36
		x := centerX + math.Cos(angle)*radius
		y := centerY + math.Sin(angle)*radius
		polygonB = append(polygonB, &draw.Point{X: x, Y: y})
	}

	// 1. 测试叠加效果
	opOverlay := polygonops.NewPolygonOverlay(polygonA, polygonB).
		WithDrawOutline(true).
		WithOutlineWidth(1.5)

	absPath, err := polygonops.DrawPolygonsWithOperations(
		exampleImageURL,
		[]*polygonops.PolygonOperation{opOverlay},
		"polygon_operations",
		"basic_overlay.png",
	)

	if err != nil {
		t.Errorf("绘制多边形叠加效果错误: %v", err)
	} else {
		fmt.Printf("成功生成多边形叠加效果图片，路径: %s\n", absPath)
	}

	// 2. 测试差集 A-B（矩形减去圆形）
	opDiffAB := polygonops.NewPolygonDifferenceAB(polygonA, polygonB).
		WithFillColor(polygonops.ColorRed).
		WithOutlineWidth(1.5)

	absPath2, err := polygonops.DrawPolygonsWithOperations(
		exampleImageURL,
		[]*polygonops.PolygonOperation{opDiffAB},
		"polygon_operations",
		"basic_diff_A_B.png",
	)

	if err != nil {
		t.Errorf("绘制多边形差集A-B错误: %v", err)
	} else {
		fmt.Printf("成功生成多边形差集A-B图片，路径: %s\n", absPath2)
	}

	// 3. 测试差集 B-A（圆形减去矩形）
	opDiffBA := polygonops.NewPolygonDifferenceBA(polygonA, polygonB).
		WithFillColor(polygonops.ColorGreen).
		WithOutlineWidth(1.5)

	absPath3, err := polygonops.DrawPolygonsWithOperations(
		exampleImageURL,
		[]*polygonops.PolygonOperation{opDiffBA},
		"polygon_operations",
		"basic_diff_B_A.png",
	)

	if err != nil {
		t.Errorf("绘制多边形差集B-A错误: %v", err)
	} else {
		fmt.Printf("成功生成多边形差集B-A图片，路径: %s\n", absPath3)
	}

	// 4. 测试交集（矩形与圆形的交集）
	opIntersection := polygonops.NewPolygonIntersection(polygonA, polygonB).
		WithFillColor(polygonops.Color{R: 0.0, G: 0.0, B: 1.0, A: 0.6}). // 蓝色
		WithOutlineWidth(1.5)

	absPath4, err := polygonops.DrawPolygonsWithOperations(
		exampleImageURL,
		[]*polygonops.PolygonOperation{opIntersection},
		"polygon_operations",
		"basic_intersection.png",
	)

	if err != nil {
		t.Errorf("绘制多边形交集错误: %v", err)
	} else {
		fmt.Printf("成功生成多边形交集图片，路径: %s\n", absPath4)
	}

	// 5. 测试组合效果（在一张图上显示三种操作）
	combinedOps := []*polygonops.PolygonOperation{
		polygonops.NewPolygonDifferenceAB(polygonA, polygonB).WithFillColor(polygonops.Color{R: 1.0, G: 0.0, B: 0.0, A: 0.4}),
		polygonops.NewPolygonDifferenceBA(polygonA, polygonB).WithFillColor(polygonops.Color{R: 0.0, G: 1.0, B: 0.0, A: 0.4}),
		polygonops.NewPolygonIntersection(polygonA, polygonB).WithFillColor(polygonops.Color{R: 0.0, G: 0.0, B: 1.0, A: 0.4}),
	}

	absPath5, err := polygonops.DrawPolygonsWithOperations(
		exampleImageURL,
		combinedOps,
		"polygon_operations",
		"basic_combined.png",
	)

	if err != nil {
		t.Errorf("绘制组合多边形操作错误: %v", err)
	} else {
		fmt.Printf("成功生成组合多边形操作图片，路径: %s\n", absPath5)
	}
}

// TestHollowPolygon 测试镂空多边形（使用even-odd填充规则）
func TestHollowPolygon(t *testing.T) {
	fmt.Println("测试镂空多边形效果（外部多边形减去内部多边形）...")

	// 创建外部多边形A的点（矩形）
	outerPoints := []*draw.Point{
		{X: 100, Y: 100}, // 左上
		{X: 500, Y: 100}, // 右上
		{X: 500, Y: 400}, // 右下
		{X: 100, Y: 400}, // 左下
	}

	// 创建内部多边形B的点（矩形）
	innerPoints := []*draw.Point{
		{X: 200, Y: 150}, // 左上
		{X: 400, Y: 150}, // 右上
		{X: 400, Y: 350}, // 右下
		{X: 200, Y: 350}, // 左下
	}

	// 创建处理器 - 使用even-odd填充规则实现镂空效果
	processor := draw.NewImageProcessor(
		exampleImageURL,
		draw.WithOutputDir("polygon_operations"),
		draw.WithOutputName("hollow_simple.png"),
		draw.WithPreProcess(func(dc *gg.Context, width, height float64) error {
			// 设置灰色半透明颜色
			dc.SetRGBA(0.5, 0.5, 0.5, 0.5)

			// 开始绘制外部多边形
			dc.MoveTo(outerPoints[0].X, outerPoints[0].Y)
			for i := 1; i < len(outerPoints); i++ {
				dc.LineTo(outerPoints[i].X, outerPoints[i].Y)
			}
			dc.ClosePath()

			// 创建内部多边形路径（注意不要填充）
			dc.NewSubPath()
			dc.MoveTo(innerPoints[0].X, innerPoints[0].Y)
			for i := 1; i < len(innerPoints); i++ {
				dc.LineTo(innerPoints[i].X, innerPoints[i].Y)
			}
			dc.ClosePath()

			// 使用even-odd填充规则，确保内部区域不填充
			dc.SetFillRule(gg.FillRuleEvenOdd)
			dc.Fill()

			// 绘制轮廓
			dc.SetLineWidth(2.0)
			dc.SetRGB(0, 0, 0)

			// 绘制外部多边形轮廓
			dc.MoveTo(outerPoints[0].X, outerPoints[0].Y)
			for i := 1; i < len(outerPoints); i++ {
				dc.LineTo(outerPoints[i].X, outerPoints[i].Y)
			}
			dc.ClosePath()
			dc.Stroke()

			// 绘制内部多边形轮廓
			dc.MoveTo(innerPoints[0].X, innerPoints[0].Y)
			for i := 1; i < len(innerPoints); i++ {
				dc.LineTo(innerPoints[i].X, innerPoints[i].Y)
			}
			dc.ClosePath()
			dc.Stroke()

			return nil
		}),
	)

	// 处理图像并保存
	outputPath, err := processor.Process()
	if err != nil {
		t.Fatalf("处理镂空多边形图像失败: %v", err)
	}

	fmt.Printf("镂空多边形图像已保存至: %s\n", outputPath)

	// 使用多边形布尔运算API实现相同效果
	opDiffAB := polygonops.NewPolygonDifferenceAB(outerPoints, innerPoints).
		WithFillColor(polygonops.Color{R: 0.5, G: 0.5, B: 0.5, A: 0.5}).
		WithDrawOutline(true).
		WithOutlineWidth(2.0)

	absPath2, err := polygonops.DrawPolygonsWithOperations(
		exampleImageURL,
		[]*polygonops.PolygonOperation{opDiffAB},
		"polygon_operations",
		"hollow_api.png",
	)

	if err != nil {
		t.Errorf("使用API绘制镂空多边形错误: %v", err)
	} else {
		fmt.Printf("成功使用API生成镂空多边形图片，路径: %s\n", absPath2)
	}
}

// TestAdvancedPolygonShapes 测试复杂多边形形状
func TestAdvancedPolygonShapes(t *testing.T) {
	fmt.Println("测试复杂多边形形状和布尔运算...")

	// 创建一个圆角矩形（多边形A）
	polygonA := []*draw.Point{}
	rectX1, rectY1 := 150.0, 150.0
	rectX2, rectY2 := 350.0, 350.0
	cornerRadius := 30.0

	// 添加左上角的点 (用8个点近似四分之一圆)
	for i := 0; i < 9; i++ {
		angle := math.Pi/2 + float64(i)*math.Pi/4/8
		x := rectX1 + cornerRadius - math.Cos(angle)*cornerRadius
		y := rectY1 + cornerRadius - math.Sin(angle)*cornerRadius
		polygonA = append(polygonA, &draw.Point{X: x, Y: y})
	}

	// 添加右上角的点
	for i := 0; i < 9; i++ {
		angle := math.Pi + float64(i)*math.Pi/4/8
		x := rectX2 - cornerRadius - math.Cos(angle)*cornerRadius
		y := rectY1 + cornerRadius - math.Sin(angle)*cornerRadius
		polygonA = append(polygonA, &draw.Point{X: x, Y: y})
	}

	// 添加右下角的点
	for i := 0; i < 9; i++ {
		angle := math.Pi*3/2 + float64(i)*math.Pi/4/8
		x := rectX2 - cornerRadius - math.Cos(angle)*cornerRadius
		y := rectY2 - cornerRadius - math.Sin(angle)*cornerRadius
		polygonA = append(polygonA, &draw.Point{X: x, Y: y})
	}

	// 添加左下角的点
	for i := 0; i < 9; i++ {
		angle := 0 + float64(i)*math.Pi/4/8
		x := rectX1 + cornerRadius - math.Cos(angle)*cornerRadius
		y := rectY2 - cornerRadius - math.Sin(angle)*cornerRadius
		polygonA = append(polygonA, &draw.Point{X: x, Y: y})
	}

	// 创建一个椭圆（多边形B）
	polygonB := []*draw.Point{}
	centerX, centerY := 250.0, 250.0
	radiusX, radiusY := 120.0, 80.0

	for i := 0; i < 36; i++ {
		angle := float64(i) * 2 * math.Pi / 36
		x := centerX + math.Cos(angle)*radiusX
		y := centerY + math.Sin(angle)*radiusY
		polygonB = append(polygonB, &draw.Point{X: x, Y: y})
	}

	// 显示所有布尔运算结果在一个图像中
	combinedOps := []*polygonops.PolygonOperation{
		polygonops.NewPolygonDifferenceAB(polygonA, polygonB).
			WithFillColor(polygonops.Color{R: 1.0, G: 0.0, B: 0.0, A: 0.4}),
		polygonops.NewPolygonDifferenceBA(polygonA, polygonB).
			WithFillColor(polygonops.Color{R: 0.0, G: 1.0, B: 0.0, A: 0.4}),
		polygonops.NewPolygonIntersection(polygonA, polygonB).
			WithFillColor(polygonops.Color{R: 0.0, G: 0.0, B: 1.0, A: 0.4}),
	}

	absPath, err := polygonops.DrawPolygonsWithOperations(
		exampleImageURL,
		combinedOps,
		"polygon_operations",
		"advanced_shapes.png",
	)

	if err != nil {
		t.Errorf("绘制高级多边形形状错误: %v", err)
	} else {
		fmt.Printf("成功生成高级多边形形状图片，路径: %s\n", absPath)
	}
}

// TestComplexHollowPolygon 测试创建复杂镂空多边形
func TestComplexHollowPolygon(t *testing.T) {
	fmt.Println("测试复杂镂空多边形效果...")

	// 创建外部多边形A的点（六边形）
	outerPoints := []*draw.Point{
		{X: 300, Y: 100}, // 顶部中点
		{X: 500, Y: 200}, // 右上
		{X: 500, Y: 400}, // 右下
		{X: 300, Y: 500}, // 底部中点
		{X: 100, Y: 400}, // 左下
		{X: 100, Y: 200}, // 左上
	}

	// 创建内部多边形B的点（五角星）
	innerPoints := []*draw.Point{
		{X: 300, Y: 150}, // 顶部点
		{X: 350, Y: 250}, // 右上点
		{X: 450, Y: 250}, // 右点
		{X: 375, Y: 325}, // 右下点
		{X: 400, Y: 425}, // 底部右点
		{X: 300, Y: 375}, // 底部中点
		{X: 200, Y: 425}, // 底部左点
		{X: 225, Y: 325}, // 左下点
		{X: 150, Y: 250}, // 左点
		{X: 250, Y: 250}, // 左上点
	}

	// 创建处理器
	processor := draw.NewImageProcessor(
		exampleImageURL,
		draw.WithOutputDir("polygon_operations"),
		draw.WithOutputName("complex_hollow.png"),
		draw.WithPreProcess(func(dc *gg.Context, width, height float64) error {
			// 设置灰色半透明颜色
			dc.SetRGBA(0.5, 0.5, 0.5, 0.7)

			// 开始绘制外部多边形
			dc.MoveTo(outerPoints[0].X, outerPoints[0].Y)
			for i := 1; i < len(outerPoints); i++ {
				dc.LineTo(outerPoints[i].X, outerPoints[i].Y)
			}
			dc.ClosePath()

			// 创建内部多边形路径（注意不要填充）
			dc.NewSubPath()
			dc.MoveTo(innerPoints[0].X, innerPoints[0].Y)
			for i := 1; i < len(innerPoints); i++ {
				dc.LineTo(innerPoints[i].X, innerPoints[i].Y)
			}
			dc.ClosePath()

			// 使用even-odd填充规则，确保内部区域不填充
			dc.SetFillRule(gg.FillRuleEvenOdd)
			dc.Fill()

			// 绘制轮廓
			dc.SetLineWidth(3.0)
			dc.SetRGB(0, 0, 0)

			// 绘制外部多边形轮廓
			dc.MoveTo(outerPoints[0].X, outerPoints[0].Y)
			for i := 1; i < len(outerPoints); i++ {
				dc.LineTo(outerPoints[i].X, outerPoints[i].Y)
			}
			dc.ClosePath()
			dc.Stroke()

			// 绘制内部多边形轮廓
			dc.MoveTo(innerPoints[0].X, innerPoints[0].Y)
			for i := 1; i < len(innerPoints); i++ {
				dc.LineTo(innerPoints[i].X, innerPoints[i].Y)
			}
			dc.ClosePath()
			dc.Stroke()

			return nil
		}),
	)

	// 处理图像并保存
	outputPath, err := processor.Process()
	if err != nil {
		t.Fatalf("处理复杂镂空多边形图像失败: %v", err)
	}

	fmt.Printf("复杂镂空多边形图像已保存至: %s\n", outputPath)

	// 使用多边形布尔运算API实现相同效果
	opDiffAB := polygonops.NewPolygonDifferenceAB(outerPoints, innerPoints).
		WithFillColor(polygonops.Color{R: 0.5, G: 0.5, B: 0.5, A: 0.7}).
		WithDrawOutline(true).
		WithOutlineWidth(3.0)

	absPath2, err := polygonops.DrawPolygonsWithOperations(
		exampleImageURL,
		[]*polygonops.PolygonOperation{opDiffAB},
		"polygon_operations",
		"complex_hollow_api.png",
	)

	if err != nil {
		t.Errorf("使用API绘制复杂镂空多边形错误: %v", err)
	} else {
		fmt.Printf("成功使用API生成复杂镂空多边形图片，路径: %s\n", absPath2)
	}
}
