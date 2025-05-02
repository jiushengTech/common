package shape

import (
	"fmt"
	"sort"

	"github.com/fogleman/gg"
	"golang.org/x/image/font/basicfont"
)

// LineType 线条类型
type LineType string

// 支持的线条类型
const (
	VerticalLine   LineType = "vertical"   // 竖线
	HorizontalLine LineType = "horizontal" // 横线
)

// Line 表示线条图形
type Line struct {
	BaseShape
	Type   LineType  `json:"line_type"` // 线条类型
	Values []float64 `json:"values"`    // 点之间的值（长度比点少1）
}

// NewLine 创建一条线
func NewLine(lineType LineType, points []Point, values []float64, options ...Option) *Line {
	line := &Line{
		BaseShape: BaseShape{
			ShapeType: "line",
			Points:    points,
			Color:     ColorYellow, // 默认黄色
			LineWidth: 2.0,
		},
		Type:   lineType,
		Values: values,
	}

	// 应用所有选项
	for _, option := range options {
		option(line)
	}

	return line
}

// Draw 实现Shape接口的绘制方法
func (l *Line) Draw(dc *gg.Context, width, height float64) error {
	switch l.Type {
	case VerticalLine:
		return l.drawVertical(dc, width, height)
	case HorizontalLine:
		return l.drawHorizontal(dc, width, height)
	default:
		return fmt.Errorf("不支持的线条类型: %s", l.Type)
	}
}

// drawVertical 绘制垂直线
func (l *Line) drawVertical(dc *gg.Context, width, height float64) error {
	for i, point := range l.Points {
		// 验证坐标
		if point.X < 0 || point.X > width {
			return fmt.Errorf("x坐标 %.2f 超出范围 [0, %.2f]", point.X, width)
		}

		// 画线
		dc.SetRGB(l.Color[0], l.Color[1], l.Color[2])
		dc.SetLineWidth(l.LineWidth)
		dc.DrawLine(point.X, 0, point.X, height)
		dc.Stroke()

		// 如果不是最后一个点，绘制值
		if i < len(l.Points)-1 && i < len(l.Values) {
			// 计算文字位置（两条线的中间位置）
			textX := (point.X + l.Points[i+1].X) / 2
			textY := height / 3 // 在1/3处显示文字

			// 绘制文字
			text := fmt.Sprintf("%.2f", l.Values[i])
			drawText(dc, text, textX, textY)
		}
	}
	return nil
}

// drawHorizontal 绘制水平线
func (l *Line) drawHorizontal(dc *gg.Context, width, height float64) error {
	for i, point := range l.Points {
		// 验证坐标
		if point.Y < 0 || point.Y > height {
			return fmt.Errorf("y坐标 %.2f 超出范围 [0, %.2f]", point.Y, height)
		}

		// 画线
		dc.SetRGB(l.Color[0], l.Color[1], l.Color[2])
		dc.SetLineWidth(l.LineWidth)
		dc.DrawLine(0, point.Y, width, point.Y)
		dc.Stroke()

		// 如果不是最后一个点，绘制值
		if i < len(l.Points)-1 && i < len(l.Values) {
			// 计算文字位置（两条线的中间位置）
			textY := (point.Y + l.Points[i+1].Y) / 2
			textX := width / 3 // 在1/3处显示文字

			// 绘制文字
			text := fmt.Sprintf("%.2f", l.Values[i])
			drawText(dc, text, textX, textY)
		}
	}
	return nil
}

// drawText 绘制文本（带描边效果）
func drawText(dc *gg.Context, text string, x, y float64) {
	// 设置文字
	face := basicfont.Face7x13
	dc.SetFontFace(face)

	// 绘制黑色描边
	dc.SetLineWidth(3)
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(text, x, y, 0.5, 0.5)
	dc.Stroke()

	// 绘制白色文字
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(text, x, y, 0.5, 0.5)
	dc.Fill()
}

// LineFactory 创建各种线条的工厂
type LineFactory struct {
	LineType LineType
}

// Create 创建线条
func (f LineFactory) Create(options ...Option) Shape {
	line := &Line{
		BaseShape: BaseShape{
			ShapeType: "line",
			Points:    []Point{},
			Color:     ColorYellow, // 默认黄色
			LineWidth: 2.0,
		},
		Type:   f.LineType,
		Values: []float64{},
	}

	// 应用所有选项
	for _, option := range options {
		option(line)
	}

	// 根据线条类型对点进行排序
	if len(line.Points) > 1 {
		if line.Type == VerticalLine {
			// 按X坐标排序
			sort.Slice(line.Points, func(i, j int) bool {
				return line.Points[i].X < line.Points[j].X
			})
		} else if line.Type == HorizontalLine {
			// 按Y坐标排序
			sort.Slice(line.Points, func(i, j int) bool {
				return line.Points[i].Y < line.Points[j].Y
			})
		}
	}

	return line
}

// WithValues 设置线条的值集合
func WithValues(values []float64) Option {
	return func(s interface{}) {
		if line, ok := s.(*Line); ok {
			line.Values = values
		}
	}
}

// WithLineType 设置线条类型
func WithLineType(lineType LineType) Option {
	return func(s interface{}) {
		if line, ok := s.(*Line); ok {
			line.Type = lineType
		}
	}
}
