package polygon

import (
	"fmt"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/jiushengTech/common/utils/draw/shape/base"
)

// Polygon 表示多边形图形
type Polygon struct {
	base.BaseShape
	Fill bool `json:"fill"` // 是否填充
}

// SetColor 设置颜色
func (p *Polygon) SetColor(color *color.RGBA) {
	p.Color = color
}

// SetLineWidth 设置线宽
func (p *Polygon) SetLineWidth(width float64) {
	p.LineWidth = width
}

// SetPoints 设置点集合
func (p *Polygon) SetPoints(points []*base.Point) {
	p.Points = points
}

// New 创建一个新的多边形
// points 需要包含多边形的所有顶点，至少需要3个点
func New(points []*base.Point, options ...base.Option) *Polygon {
	polygon := &Polygon{
		BaseShape: base.BaseShape{
			ShapeType: "polygon",
			Points:    points,
			Color:     base.ColorGreen, // 默认绿色
			LineWidth: 2.0,
		},
		Fill: false,
	}

	// 应用所有选项
	for _, option := range options {
		option(polygon)
	}

	return polygon
}

// Draw 实现Shape接口的绘制方法
func (p *Polygon) Draw(dc *gg.Context, width, height float64) error {
	if len(p.Points) < 3 {
		return fmt.Errorf("多边形至少需要3个点")
	}

	// 设置颜色和线宽
	dc.SetColor(p.Color)
	dc.SetLineWidth(p.LineWidth)

	// 开始绘制路径
	dc.MoveTo(p.Points[0].X, p.Points[0].Y)

	// 添加所有点到路径
	for i := 1; i < len(p.Points); i++ {
		dc.LineTo(p.Points[i].X, p.Points[i].Y)
	}

	// 闭合路径
	dc.ClosePath()

	// 绘制
	if p.Fill {
		dc.Fill()
	} else {
		dc.Stroke()
	}

	return nil
}

// Factory 创建多边形的工厂
type Factory struct{}

// Create 创建多边形
func (f Factory) Create(options ...base.Option) base.Shape {
	polygon := &Polygon{
		BaseShape: base.BaseShape{
			ShapeType: "polygon",
			Points:    []*base.Point{},
			Color:     base.ColorGreen, // 默认绿色
			LineWidth: 2.0,
		},
		Fill: false,
	}

	// 应用所有选项
	for _, option := range options {
		option(polygon)
	}

	return polygon
}

// WithFill 设置是否填充
func WithFill(fill bool) base.Option {
	return func(s interface{}) {
		if polygon, ok := s.(*Polygon); ok {
			polygon.Fill = fill
		}
	}
}
