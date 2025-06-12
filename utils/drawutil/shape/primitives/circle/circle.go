package circle

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/jiushengTech/common/utils/drawutil/colorx"
	"github.com/jiushengTech/common/utils/drawutil/shape/base"
)

// Circle 表示圆形图形
type Circle struct {
	base.BaseShape
	Radius float64 `json:"radius"` // 圆的半径
	Fill   bool    `json:"fill"`   // 是否填充
}

// New 创建一个新的圆形
// radius 圆的半径，fill 是否填充
func New(radius float64, fill bool, options ...base.Option) *Circle {
	circle := &Circle{
		BaseShape: base.BaseShape{
			ShapeType: "circle",
			Points:    []*base.Point{},
			Color:     colorx.Green,
			LineWidth: 2.0,
		},
		Radius: radius,
		Fill:   fill,
	}

	// 应用所有选项
	base.ApplyOptions(&circle.BaseShape, options...)

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
	dc.SetColor(c.Color)
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
