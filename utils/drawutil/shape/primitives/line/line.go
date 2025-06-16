package line

import (
	"fmt"
	"sort"

	"github.com/fogleman/gg"
	"github.com/jiushengTech/common/utils/drawutil/colorx"
	"github.com/jiushengTech/common/utils/drawutil/shape/base"
	"golang.org/x/image/font/basicfont"
)

// Type 线条类型
type Type string

// 支持的线条类型
const (
	Vertical   Type = "vertical"   // 竖线
	Horizontal Type = "horizontal" // 横线
)

// Line 表示线条图形
type Line struct {
	base.BaseShape
	LineOptions
}

// New 创建一条线
func New(options ...LineOption) *Line {
	// 默认配置
	lineOpts := LineOptions{
		Type:         Vertical,
		Values:       []float64{},
		TextPosition: 0.5,
	}

	line := &Line{
		BaseShape: base.BaseShape{
			ShapeType: "line",
			Points:    []*base.Point{},
			Color:     colorx.Yellow,
			LineWidth: 2.0,
		},
		LineOptions: lineOpts,
	}

	// 应用所有选项
	for _, opt := range options {
		opt(&line.LineOptions, &line.BaseShape)
	}

	// 根据线条类型对点进行排序
	line.sortPoints()

	return line
}

// sortPoints 根据线条类型对点进行排序
func (l *Line) sortPoints() {
	if len(l.Points) > 1 {
		if l.Type == Vertical {
			// 按X坐标排序
			sort.Slice(l.Points, func(i, j int) bool {
				return l.Points[i].X < l.Points[j].X
			})
		} else if l.Type == Horizontal {
			// 按Y坐标排序
			sort.Slice(l.Points, func(i, j int) bool {
				return l.Points[i].Y < l.Points[j].Y
			})
		}
	}
}

// Draw 实现Shape接口的绘制方法
func (l *Line) Draw(dc *gg.Context, width, height float64) error {
	switch l.Type {
	case Vertical:
		return l.drawVertical(dc, width, height)
	case Horizontal:
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
		dc.SetColor(l.Color)
		dc.SetLineWidth(l.LineWidth)
		dc.DrawLine(point.X, 0, point.X, height)
		dc.Stroke()

		// 如果不是最后一个点，绘制值
		if i < len(l.Points)-1 && i < len(l.Values) {
			// 计算文字位置（两条线的中间位置）
			textX := (point.X + l.Points[i+1].X) / 2
			textY := height * l.TextPosition // 使用设置的位置比例

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
		dc.SetColor(l.Color)
		dc.SetLineWidth(l.LineWidth)
		dc.DrawLine(0, point.Y, width, point.Y)
		dc.Stroke()

		// 如果不是最后一个点，绘制值
		if i < len(l.Points)-1 && i < len(l.Values) {
			// 计算文字位置（两条线的中间位置）
			textY := (point.Y + l.Points[i+1].Y) / 2
			textX := width * l.TextPosition // 使用设置的位置比例

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
