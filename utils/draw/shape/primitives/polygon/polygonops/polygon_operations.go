package polygonops

import (
	"fmt"
	"github.com/jiushengTech/common/utils/draw/colorx"
	"github.com/jiushengTech/common/utils/draw/processor"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/jiushengTech/common/utils/draw/shape/base"
)

// PolygonOperation 表示一个多边形操作
type PolygonOperation struct {
	// 操作类型
	Type OperationType
	// 要处理的多边形A
	PolygonA []*base.Point
	// 要处理的多边形B
	PolygonB []*base.Point
	// 填充颜色 (RGBA)
	FillColor *color.RGBA
	// 是否绘制轮廓
	DrawOutline bool
	// 轮廓颜色
	OutlineColor *color.RGBA
	// 轮廓宽度
	OutlineWidth float64
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

// NewPolygonOverlay 创建两个叠加显示的多边形操作
func NewPolygonOverlay(polygonA, polygonB []*base.Point, options ...Option) *PolygonOperation {
	po := &PolygonOperation{
		Type:         OperationOverlay,
		PolygonA:     polygonA,
		PolygonB:     polygonB,
		DrawOutline:  false,
		OutlineWidth: 1.0,
	}

	// 应用选项
	for _, option := range options {
		option(po)
	}

	return po
}

// NewPolygonDifferenceAB 创建多边形差集操作 A-B
func NewPolygonDifferenceAB(polygonA, polygonB []*base.Point, options ...Option) *PolygonOperation {
	po := &PolygonOperation{
		Type:         OperationDifferenceAB,
		PolygonA:     polygonA,
		PolygonB:     polygonB,
		DrawOutline:  true,
		OutlineWidth: 1.0,
	}

	// 应用选项
	for _, option := range options {
		option(po)
	}

	return po
}

// NewPolygonDifferenceBA 创建多边形差集操作 B-A
func NewPolygonDifferenceBA(polygonA, polygonB []*base.Point, options ...Option) *PolygonOperation {
	po := &PolygonOperation{
		Type:         OperationDifferenceBA,
		PolygonA:     polygonA,
		PolygonB:     polygonB,
		DrawOutline:  true,
		OutlineWidth: 1.0,
	}

	// 应用选项
	for _, option := range options {
		option(po)
	}

	return po
}

// NewPolygonIntersection 创建多边形交集操作
func NewPolygonIntersection(polygonA, polygonB []*base.Point, options ...Option) *PolygonOperation {
	po := &PolygonOperation{
		Type:         OperationIntersection,
		PolygonA:     polygonA,
		PolygonB:     polygonB,
		DrawOutline:  true,
		OutlineWidth: 1.0,
	}

	// 应用选项
	for _, option := range options {
		option(po)
	}

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
		dc.SetColor(po.FillColor)

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
		dc.SetColor(colorx.Green)

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
	dc.SetColor(po.FillColor)

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
	dc.SetColor(po.FillColor)

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

// drawIntersection 绘制交集 A∩B
func (po *PolygonOperation) drawIntersection(dc *gg.Context) error {
	// 使用剪裁技术实现交集
	dc.SetColor(po.FillColor)

	// 创建路径A
	if len(po.PolygonA) > 0 {
		dc.MoveTo(po.PolygonA[0].X, po.PolygonA[0].Y)
		for i := 1; i < len(po.PolygonA); i++ {
			dc.LineTo(po.PolygonA[i].X, po.PolygonA[i].Y)
		}
		dc.ClosePath()
	}

	// 设置为剪裁路径
	dc.Clip()

	// 绘制路径B，它会被路径A剪裁
	if len(po.PolygonB) > 0 {
		dc.MoveTo(po.PolygonB[0].X, po.PolygonB[0].Y)
		for i := 1; i < len(po.PolygonB); i++ {
			dc.LineTo(po.PolygonB[i].X, po.PolygonB[i].Y)
		}
		dc.ClosePath()
	}

	// 填充交集区域
	dc.Fill()

	// 重置剪裁
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
	dc.SetColor(po.OutlineColor)
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
// 已废弃: 请使用 DrawPolygons 函数代替
func DrawPolygonsWithOperations(imageURL string, operations []*PolygonOperation, outputDir, outputName string) (string, error) {
	return DrawPolygons(
		WithImageURL(imageURL),
		WithOperations(operations),
		WithOutputDirectory(outputDir),
		WithOutputFileName(outputName),
	)
}

// DrawPolygons 使用选项模式绘制多边形操作
func DrawPolygons(options ...DrawOption) (string, error) {
	// 创建默认配置
	cfg := &DrawConfig{
		OutputDir:  "polygon_output",                                    // 默认输出目录
		OutputName: processor.GetDefaultOutputName(processor.FormatPNG), // 默认输出文件名
	}

	// 应用选项
	for _, option := range options {
		option(cfg)
	}

	// 验证必要的配置
	if cfg.ImageURL == "" {
		return "", fmt.Errorf("image URL is required")
	}

	if len(cfg.Operations) == 0 {
		return "", fmt.Errorf("at least one polygon operation is required")
	}

	// 创建处理器
	p := processor.NewImageProcessor(
		cfg.ImageURL,
		processor.WithOutputName(cfg.OutputName),
		processor.WithOutputDir(cfg.OutputDir),
		processor.WithPreProcess(func(dc *gg.Context, width, height float64) error {
			// 依次执行所有多边形操作
			for _, op := range cfg.Operations {
				if err := op.Draw(dc); err != nil {
					return err
				}
			}
			return nil
		}),
	)

	// 处理图像并返回路径
	return p.Process()
}
