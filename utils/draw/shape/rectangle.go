package shape

import (
	"fmt"

	"github.com/fogleman/gg"
)

// Rectangle 表示矩形图形
type Rectangle struct {
	BaseShape
	Fill bool `json:"fill"` // 是否填充
}

// NewRectangle 创建一个新的矩形
// points 需要包含对角线的两个点: [左上角, 右下角]
func NewRectangle(points []Point, options ...Option) *Rectangle {
	rect := &Rectangle{
		BaseShape: BaseShape{
			ShapeType: "rectangle",
			Points:    points,
			Color:     ColorYellow, // 默认黄色
			LineWidth: 2.0,
		},
		Fill: false,
	}

	// 应用所有选项
	for _, option := range options {
		option(rect)
	}

	return rect
}

// Draw 实现Shape接口的绘制方法
func (r *Rectangle) Draw(dc *gg.Context, width, height float64) error {
	if len(r.Points) < 2 {
		return fmt.Errorf("矩形需要至少两个点")
	}

	// 获取左上角和右下角点
	topLeft := r.Points[0]
	bottomRight := r.Points[1]

	// 验证坐标
	if topLeft.X < 0 || topLeft.X > width || topLeft.Y < 0 || topLeft.Y > height ||
		bottomRight.X < 0 || bottomRight.X > width || bottomRight.Y < 0 || bottomRight.Y > height {
		return fmt.Errorf("矩形坐标超出范围")
	}

	// 计算宽度和高度
	rectWidth := bottomRight.X - topLeft.X
	rectHeight := bottomRight.Y - topLeft.Y

	// 设置颜色和线宽
	dc.SetRGB(r.Color[0], r.Color[1], r.Color[2])
	dc.SetLineWidth(r.LineWidth)

	// 绘制矩形
	if r.Fill {
		dc.DrawRectangle(topLeft.X, topLeft.Y, rectWidth, rectHeight)
		dc.Fill()
	} else {
		dc.DrawRectangle(topLeft.X, topLeft.Y, rectWidth, rectHeight)
		dc.Stroke()
	}

	return nil
}

// RectangleFactory 创建矩形的工厂
type RectangleFactory struct{}

// Create 创建矩形
func (f RectangleFactory) Create(options ...Option) Shape {
	rect := &Rectangle{
		BaseShape: BaseShape{
			ShapeType: "rectangle",
			Points:    []Point{},
			Color:     ColorYellow, // 默认黄色
			LineWidth: 2.0,
		},
		Fill: false,
	}

	// 应用所有选项
	for _, option := range options {
		option(rect)
	}

	return rect
}

// WithFill 设置是否填充矩形
func WithFill(fill bool) Option {
	return func(s interface{}) {
		if rect, ok := s.(*Rectangle); ok {
			rect.Fill = fill
		}
	}
}
