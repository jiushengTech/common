package circle

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/jiushengTech/common/utils/draw/shape/base"
)

// Circle 表示圆形图形
type Circle struct {
	base.BaseShape
	Radius float64 `json:"radius"` // 圆的半径
	Fill   bool    `json:"fill"`   // 是否填充
}

// SetColor 设置颜色
func (c *Circle) SetColor(color [3]float64) {
	c.Color = color
}

// SetLineWidth 设置线宽
func (c *Circle) SetLineWidth(width float64) {
	c.LineWidth = width
}

// SetPoints 设置点集合
func (c *Circle) SetPoints(points []*base.Point) {
	c.Points = points
}

// New 创建一个新的圆形
// points 需要包含圆心
func New(center *base.Point, radius float64, options ...base.Option) *Circle {
	circle := &Circle{
		BaseShape: base.BaseShape{
			ShapeType: "circle",
			Points:    []*base.Point{center},
			Color:     base.ColorYellow, // 默认黄色
			LineWidth: 2.0,
		},
		Radius: radius,
		Fill:   false,
	}

	// 应用所有选项
	for _, option := range options {
		option(circle)
	}

	return circle
}

// Draw 实现Shape接口的绘制方法
func (c *Circle) Draw(dc *gg.Context, width, height float64) error {
	if len(c.Points) < 1 {
		return fmt.Errorf("圆形需要圆心点")
	}

	// 获取圆心
	center := c.Points[0]

	// 验证坐标
	if center.X-c.Radius < 0 || center.X+c.Radius > width ||
		center.Y-c.Radius < 0 || center.Y+c.Radius > height {
		return fmt.Errorf("圆形范围超出图像边界")
	}

	// 设置颜色和线宽
	dc.SetRGB(c.Color[0], c.Color[1], c.Color[2])
	dc.SetLineWidth(c.LineWidth)

	// 绘制圆形
	dc.DrawCircle(center.X, center.Y, c.Radius)
	if c.Fill {
		dc.Fill()
	} else {
		dc.Stroke()
	}

	return nil
}

// Factory 创建圆形的工厂
type Factory struct{}

// Create 创建圆形
func (f Factory) Create(options ...base.Option) base.Shape {
	circle := &Circle{
		BaseShape: base.BaseShape{
			ShapeType: "circle",
			Points:    []*base.Point{},
			Color:     base.ColorYellow, // 默认黄色
			LineWidth: 2.0,
		},
		Radius: 10, // 默认半径
		Fill:   false,
	}

	// 应用所有选项
	for _, option := range options {
		option(circle)
	}

	return circle
}

// WithRadius 设置圆的半径
func WithRadius(radius float64) base.Option {
	return func(s interface{}) {
		if circle, ok := s.(*Circle); ok {
			circle.Radius = radius
		}
	}
}

// WithFill 设置是否填充
func WithFill(fill bool) base.Option {
	return func(s interface{}) {
		if circle, ok := s.(*Circle); ok {
			circle.Fill = fill
		}
	}
}
