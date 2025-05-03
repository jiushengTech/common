// Package polygonops 适配器包，为主包提供多边形布尔运算功能
package polygonops

import (
	"github.com/fogleman/gg"
	"github.com/jiushengTech/common/utils/draw/shape/base"
)

// Point类型别名，用于适配
type Point = base.Point

// Color 表示RGBA颜色
type Color struct {
	R, G, B, A float64
}

// OperationType 表示多边形操作类型
type OperationType int

const (
	// OperationOverlay 表示两个多边形叠加显示（半透明）
	OperationOverlay OperationType = iota
	// OperationDifferenceAB 表示多边形差集 A-B
	OperationDifferenceAB
	// OperationDifferenceBA 表示多边形差集 B-A
	OperationDifferenceBA
	// OperationIntersection 表示多边形交集
	OperationIntersection
)

// 预定义颜色
var (
	ColorGray  = Color{0.5, 0.5, 0.5, 0.3} // 灰色（半透明）
	ColorBlue  = Color{0.0, 0.0, 1.0, 0.3} // 蓝色（半透明）
	ColorRed   = Color{1.0, 0.0, 0.0, 0.6} // 红色（半透明）
	ColorGreen = Color{0.0, 1.0, 0.0, 0.6} // 绿色（半透明）
	ColorBlack = Color{0.0, 0.0, 0.0, 1.0} // 黑色
)

// PolygonOperation 表示一个多边形操作
type PolygonOperation struct {
	// 操作类型
	Type OperationType
	// 要处理的多边形A
	PolygonA []*Point
	// 要处理的多边形B
	PolygonB []*Point
	// 填充颜色 (RGBA)
	FillColor Color
	// 是否绘制轮廓
	DrawOutline bool
	// 轮廓颜色
	OutlineColor Color
	// 轮廓宽度
	OutlineWidth float64
}

// NewPolygonOverlay 创建叠加显示的多边形操作
func NewPolygonOverlay(polygonA, polygonB []*Point) *PolygonOperation {
	return &PolygonOperation{
		Type:         OperationOverlay,
		PolygonA:     polygonA,
		PolygonB:     polygonB,
		FillColor:    ColorGray,
		DrawOutline:  false,
		OutlineColor: ColorBlack,
		OutlineWidth: 1.0,
	}
}

// NewPolygonDifferenceAB 创建差集(A-B)多边形操作
func NewPolygonDifferenceAB(polygonA, polygonB []*Point) *PolygonOperation {
	return &PolygonOperation{
		Type:         OperationDifferenceAB,
		PolygonA:     polygonA,
		PolygonB:     polygonB,
		FillColor:    ColorRed,
		DrawOutline:  true,
		OutlineColor: ColorBlack,
		OutlineWidth: 1.0,
	}
}

// NewPolygonDifferenceBA 创建差集(B-A)多边形操作
func NewPolygonDifferenceBA(polygonA, polygonB []*Point) *PolygonOperation {
	return &PolygonOperation{
		Type:         OperationDifferenceBA,
		PolygonA:     polygonA,
		PolygonB:     polygonB,
		FillColor:    ColorGreen,
		DrawOutline:  true,
		OutlineColor: ColorBlack,
		OutlineWidth: 1.0,
	}
}

// NewPolygonIntersection 创建交集多边形操作
func NewPolygonIntersection(polygonA, polygonB []*Point) *PolygonOperation {
	return &PolygonOperation{
		Type:         OperationIntersection,
		PolygonA:     polygonA,
		PolygonB:     polygonB,
		FillColor:    ColorBlue,
		DrawOutline:  true,
		OutlineColor: ColorBlack,
		OutlineWidth: 1.0,
	}
}

// WithFillColor 设置填充颜色
func (po *PolygonOperation) WithFillColor(color Color) *PolygonOperation {
	po.FillColor = color
	return po
}

// WithOutlineColor 设置轮廓颜色
func (po *PolygonOperation) WithOutlineColor(color Color) *PolygonOperation {
	po.OutlineColor = color
	return po
}

// WithOutlineWidth 设置轮廓宽度
func (po *PolygonOperation) WithOutlineWidth(width float64) *PolygonOperation {
	po.OutlineWidth = width
	return po
}

// WithDrawOutline 设置是否绘制轮廓
func (po *PolygonOperation) WithDrawOutline(draw bool) *PolygonOperation {
	po.DrawOutline = draw
	return po
}

// Draw 将多边形操作绘制到上下文中
func (po *PolygonOperation) Draw(dc *gg.Context) error {
	switch po.Type {
	case OperationOverlay:
		return po.drawOverlay(dc)
	case OperationDifferenceAB:
		return po.drawDifferenceAB(dc)
	case OperationDifferenceBA:
		return po.drawDifferenceBA(dc)
	case OperationIntersection:
		return po.drawIntersection(dc)
	default:
		return nil
	}
}

// drawOverlay 绘制两个多边形叠加效果
func (po *PolygonOperation) drawOverlay(dc *gg.Context) error {
	// 绘制多边形A
	if len(po.PolygonA) > 0 {
		// 设置颜色
		dc.SetRGBA(po.FillColor.R, po.FillColor.G, po.FillColor.B, po.FillColor.A)

		// 绘制路径
		dc.MoveTo(po.PolygonA[0].X, po.PolygonA[0].Y)
		for i := 1; i < len(po.PolygonA); i++ {
			dc.LineTo(po.PolygonA[i].X, po.PolygonA[i].Y)
		}
		dc.ClosePath()
		dc.Fill()
	}

	// 绘制多边形B
	if len(po.PolygonB) > 0 {
		// 使用蓝色半透明
		dc.SetRGBA(ColorBlue.R, ColorBlue.G, ColorBlue.B, ColorBlue.A)

		// 绘制路径
		dc.MoveTo(po.PolygonB[0].X, po.PolygonB[0].Y)
		for i := 1; i < len(po.PolygonB); i++ {
			dc.LineTo(po.PolygonB[i].X, po.PolygonB[i].Y)
		}
		dc.ClosePath()
		dc.Fill()
	}

	// 绘制轮廓（如果需要）
	if po.DrawOutline {
		po.drawOutlines(dc)
	}

	return nil
}

// drawDifferenceAB 绘制差集 A-B
func (po *PolygonOperation) drawDifferenceAB(dc *gg.Context) error {
	// 设置颜色
	dc.SetRGBA(po.FillColor.R, po.FillColor.G, po.FillColor.B, po.FillColor.A)

	// 先绘制多边形A路径
	if len(po.PolygonA) > 0 {
		dc.MoveTo(po.PolygonA[0].X, po.PolygonA[0].Y)
		for i := 1; i < len(po.PolygonA); i++ {
			dc.LineTo(po.PolygonA[i].X, po.PolygonA[i].Y)
		}
		dc.ClosePath()
	}

	// 再绘制多边形B路径（用于剪裁掉交集区域）
	dc.NewSubPath()
	if len(po.PolygonB) > 0 {
		dc.MoveTo(po.PolygonB[0].X, po.PolygonB[0].Y)
		for i := 1; i < len(po.PolygonB); i++ {
			dc.LineTo(po.PolygonB[i].X, po.PolygonB[i].Y)
		}
		dc.ClosePath()
	}

	// 使用even-odd填充规则确保只填充A-B区域
	dc.SetFillRule(gg.FillRuleEvenOdd)
	dc.Fill()

	// 绘制轮廓（如果需要）
	if po.DrawOutline {
		po.drawOutlines(dc)
	}

	return nil
}

// drawDifferenceBA 绘制差集 B-A
func (po *PolygonOperation) drawDifferenceBA(dc *gg.Context) error {
	// 设置颜色
	dc.SetRGBA(po.FillColor.R, po.FillColor.G, po.FillColor.B, po.FillColor.A)

	// 先绘制多边形B路径
	if len(po.PolygonB) > 0 {
		dc.MoveTo(po.PolygonB[0].X, po.PolygonB[0].Y)
		for i := 1; i < len(po.PolygonB); i++ {
			dc.LineTo(po.PolygonB[i].X, po.PolygonB[i].Y)
		}
		dc.ClosePath()
	}

	// 再绘制多边形A路径（用于剪裁掉交集区域）
	dc.NewSubPath()
	if len(po.PolygonA) > 0 {
		dc.MoveTo(po.PolygonA[0].X, po.PolygonA[0].Y)
		for i := 1; i < len(po.PolygonA); i++ {
			dc.LineTo(po.PolygonA[i].X, po.PolygonA[i].Y)
		}
		dc.ClosePath()
	}

	// 使用even-odd填充规则确保只填充B-A区域
	dc.SetFillRule(gg.FillRuleEvenOdd)
	dc.Fill()

	// 绘制轮廓（如果需要）
	if po.DrawOutline {
		po.drawOutlines(dc)
	}

	return nil
}

// drawIntersection 绘制交集
func (po *PolygonOperation) drawIntersection(dc *gg.Context) error {
	// 设置颜色
	dc.SetRGBA(po.FillColor.R, po.FillColor.G, po.FillColor.B, po.FillColor.A)

	// 使用裁剪区域来实现交集
	// 首先，定义多边形A作为裁剪区域
	if len(po.PolygonA) > 0 {
		dc.MoveTo(po.PolygonA[0].X, po.PolygonA[0].Y)
		for i := 1; i < len(po.PolygonA); i++ {
			dc.LineTo(po.PolygonA[i].X, po.PolygonA[i].Y)
		}
		dc.ClosePath()
		dc.Clip() // 将当前路径设为裁剪区域
	}

	// 然后绘制多边形B，只有在裁剪区域内的部分才会显示
	if len(po.PolygonB) > 0 {
		dc.MoveTo(po.PolygonB[0].X, po.PolygonB[0].Y)
		for i := 1; i < len(po.PolygonB); i++ {
			dc.LineTo(po.PolygonB[i].X, po.PolygonB[i].Y)
		}
		dc.ClosePath()
		dc.Fill() // 填充B，但只有与A相交的部分会显示
	}

	// 重置裁剪区域
	dc.ResetClip()

	// 绘制轮廓（如果需要）
	if po.DrawOutline {
		po.drawOutlines(dc)
	}

	return nil
}

// drawOutlines 绘制多边形轮廓
func (po *PolygonOperation) drawOutlines(dc *gg.Context) {
	// 设置轮廓颜色和宽度
	dc.SetRGBA(po.OutlineColor.R, po.OutlineColor.G, po.OutlineColor.B, po.OutlineColor.A)
	dc.SetLineWidth(po.OutlineWidth)

	// 绘制多边形A轮廓
	if len(po.PolygonA) > 0 {
		dc.MoveTo(po.PolygonA[0].X, po.PolygonA[0].Y)
		for i := 1; i < len(po.PolygonA); i++ {
			dc.LineTo(po.PolygonA[i].X, po.PolygonA[i].Y)
		}
		dc.ClosePath()
		dc.Stroke()
	}

	// 绘制多边形B轮廓
	if len(po.PolygonB) > 0 {
		dc.MoveTo(po.PolygonB[0].X, po.PolygonB[0].Y)
		for i := 1; i < len(po.PolygonB); i++ {
			dc.LineTo(po.PolygonB[i].X, po.PolygonB[i].Y)
		}
		dc.ClosePath()
		dc.Stroke()
	}
}

// DrawPolygonsWithOperations 使用指定的图片和多边形操作列表生成图片
func DrawPolygonsWithOperations(imageURL string, operations []*PolygonOperation, outputDir, outputName string) (string, error) {
	// 这个函数需要直接引入原始代码，因为我们不能导入draw包
	// 参考 `shape/primitives/polygonops/polygon_operations.go` 中的相同函数
	// 这需要引入draw和processor包的代码，此处仅作示例

	// 实际实现会使用预处理函数绘制所有多边形
	// 以下是简化版本
	processor := NewImageProcessor(imageURL)
	processor.SetOutputDir(outputDir)
	processor.SetOutputName(outputName)

	// 绘制所有多边形操作
	processor.SetPreProcess(func(dc *gg.Context, width, height float64) error {
		for _, op := range operations {
			if err := op.Draw(dc); err != nil {
				return err
			}
		}
		return nil
	})

	// 处理并返回路径
	return processor.Process()
}

// NewPolygonColor 创建一个新的多边形颜色
func NewPolygonColor(r, g, b, a float64) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// ImageProcessor 简化版的图像处理器
type ImageProcessor struct {
	imagePath  string
	outputDir  string
	outputName string
	preProcess func(*gg.Context, float64, float64) error
}

// NewImageProcessor 创建新的图像处理器
func NewImageProcessor(imagePath string) *ImageProcessor {
	return &ImageProcessor{
		imagePath: imagePath,
	}
}

// SetOutputDir 设置输出目录
func (p *ImageProcessor) SetOutputDir(dir string) {
	p.outputDir = dir
}

// SetOutputName 设置输出文件名
func (p *ImageProcessor) SetOutputName(name string) {
	p.outputName = name
}

// SetPreProcess 设置预处理函数
func (p *ImageProcessor) SetPreProcess(fn func(*gg.Context, float64, float64) error) {
	p.preProcess = fn
}

// Process 处理图像
func (p *ImageProcessor) Process() (string, error) {
	// 实际实现会调用原始处理器的代码
	// 此处仅做最小化实现，实际使用时需要通过主包的draw.go接口
	return "模拟的输出路径", nil
}
