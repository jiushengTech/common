package shape

import (
	"fmt"

	"github.com/fogleman/gg"
)

// Circle 表示圆形图形
type Circle struct {
	BaseShape
	Radius float64 `json:"radius"` // 圆的半径
	Fill   bool    `json:"fill"`   // 是否填充
}

// NewCircle 创建一个新的圆形
// points 需要包含圆心
func NewCircle(center Point, radius float64, options ...Option) *Circle {
	circle := &Circle{
		BaseShape: BaseShape{
			ShapeType: "circle",
			Points:    []Point{center},
			Color:     ColorYellow, // 默认黄色
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

// CircleFactory 创建圆形的工厂
type CircleFactory struct{}

// Create 创建圆形
func (f CircleFactory) Create(options ...Option) Shape {
	circle := &Circle{
		BaseShape: BaseShape{
			ShapeType: "circle",
			Points:    []Point{},
			Color:     ColorYellow, // 默认黄色
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
func WithRadius(radius float64) Option {
	return func(s interface{}) {
		if circle, ok := s.(*Circle); ok {
			circle.Radius = radius
		}
	}
}
