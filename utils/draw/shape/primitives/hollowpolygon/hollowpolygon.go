package hollowpolygon

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/jiushengTech/common/utils/draw/shape/base"
)

// HollowPolygon 表示一个带有镂空区域的多边形
// 外部多边形包含内部多边形，内部多边形区域将被镂空
type HollowPolygon struct {
	base.BaseShape
	OuterPoints []base.Point `json:"outer_points"` // 外部多边形的点
	InnerPoints []base.Point `json:"inner_points"` // 内部多边形的点
	Opacity     float64      `json:"opacity"`      // 不透明度 (0-1)
}

// SetColor 设置颜色
func (h *HollowPolygon) SetColor(color [3]float64) {
	h.Color = color
}

// SetLineWidth 设置线宽
func (h *HollowPolygon) SetLineWidth(width float64) {
	h.LineWidth = width
}

// SetPoints 设置点集合（这里点集合只用于兼容接口，实际不使用）
func (h *HollowPolygon) SetPoints(points []base.Point) {
	h.Points = points
}

// SetOpacity 设置不透明度
func (h *HollowPolygon) SetOpacity(opacity float64) {
	h.Opacity = opacity
}

// SetOuterPoints 设置外部多边形的顶点
func (h *HollowPolygon) SetOuterPoints(points []base.Point) {
	h.OuterPoints = points
}

// SetInnerPoints 设置内部多边形的顶点
func (h *HollowPolygon) SetInnerPoints(points []base.Point) {
	h.InnerPoints = points
}

// New 创建一个新的镂空多边形
// 参数:
//
//	outerPoints: 外部多边形的顶点，至少需要3个点
//	innerPoints: 内部多边形的顶点，至少需要3个点
//	options: 可选配置
func New(outerPoints, innerPoints []base.Point, options ...base.Option) *HollowPolygon {
	hollow := &HollowPolygon{
		BaseShape: base.BaseShape{
			ShapeType: "hollow_polygon",
			Color:     base.ColorBlack, // 默认黑色
			LineWidth: 1.0,
		},
		OuterPoints: outerPoints,
		InnerPoints: innerPoints,
		Opacity:     0.5, // 默认半透明
	}

	// 应用所有选项
	for _, option := range options {
		option(hollow)
	}

	return hollow
}

// Draw 实现Shape接口的绘制方法
func (h *HollowPolygon) Draw(dc *gg.Context, width, height float64) error {
	if len(h.OuterPoints) < 3 {
		return fmt.Errorf("外部多边形至少需要3个点")
	}

	if len(h.InnerPoints) < 3 {
		return fmt.Errorf("内部多边形至少需要3个点")
	}

	// 保存当前状态
	dc.Push()

	// 设置RGBA颜色带透明度
	dc.SetRGBA(h.Color[0], h.Color[1], h.Color[2], h.Opacity)
	dc.SetLineWidth(h.LineWidth)

	// 绘制外部多边形路径
	dc.NewSubPath()
	dc.MoveTo(h.OuterPoints[0].X, h.OuterPoints[0].Y)
	for i := 1; i < len(h.OuterPoints); i++ {
		dc.LineTo(h.OuterPoints[i].X, h.OuterPoints[i].Y)
	}
	dc.ClosePath()

	// 绘制内部多边形路径（作为孔洞）
	dc.NewSubPath()
	dc.MoveTo(h.InnerPoints[0].X, h.InnerPoints[0].Y)
	for i := 1; i < len(h.InnerPoints); i++ {
		dc.LineTo(h.InnerPoints[i].X, h.InnerPoints[i].Y)
	}
	dc.ClosePath()

	// 使用偶奇填充规则(even-odd)填充路径
	// 这确保内部多边形区域被"挖空"
	dc.FillPreserve()

	// 绘制边界线
	dc.SetRGB(h.Color[0], h.Color[1], h.Color[2])
	dc.Stroke()

	// 恢复之前的状态
	dc.Pop()

	return nil
}

// Factory 创建镂空多边形的工厂
type Factory struct{}

// Create 创建镂空多边形
func (f Factory) Create(options ...base.Option) base.Shape {
	hollow := &HollowPolygon{
		BaseShape: base.BaseShape{
			ShapeType: "hollow_polygon",
			Color:     base.ColorBlack, // 默认黑色
			LineWidth: 1.0,
		},
		Opacity: 0.5, // 默认半透明
	}

	// 应用所有选项
	for _, option := range options {
		option(hollow)
	}

	return hollow
}

// WithOpacity 设置不透明度
func WithOpacity(opacity float64) base.Option {
	return func(s interface{}) {
		if h, ok := s.(*HollowPolygon); ok {
			h.Opacity = opacity
		}
	}
}

// WithOuterPoints 设置外部多边形的顶点
func WithOuterPoints(points []base.Point) base.Option {
	return func(s interface{}) {
		if h, ok := s.(*HollowPolygon); ok {
			h.OuterPoints = points
		}
	}
}

// WithInnerPoints 设置内部多边形的顶点
func WithInnerPoints(points []base.Point) base.Option {
	return func(s interface{}) {
		if h, ok := s.(*HollowPolygon); ok {
			h.InnerPoints = points
		}
	}
}
