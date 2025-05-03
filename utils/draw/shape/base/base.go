package base

import (
	"github.com/fogleman/gg"
)

// Point 表示二维坐标点
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Shape 是所有图形的通用接口
type Shape interface {
	// Draw 在给定的画布上绘制图形
	Draw(dc *gg.Context, width, height float64) error

	// GetType 返回图形的类型
	GetType() string

	// GetColor 返回图形的颜色
	GetColor() [3]float64

	// GetPoints 返回图形的点集合
	GetPoints() []*Point
}

// 颜色常量 (RGB值，范围0-1)
var (
	ColorWhite   = [3]float64{1, 1, 1}       // 白色
	ColorBlack   = [3]float64{0, 0, 0}       // 黑色
	ColorRed     = [3]float64{1, 0, 0}       // 红色
	ColorBlue    = [3]float64{0, 0, 1}       // 蓝色
	ColorGreen   = [3]float64{0, 1, 0}       // 绿色
	ColorYellow  = [3]float64{1, 1, 0}       // 黄色
	ColorCyan    = [3]float64{0, 1, 1}       // 青色
	ColorMagenta = [3]float64{1, 0, 1}       // 品红
	ColorGray    = [3]float64{0.5, 0.5, 0.5} // 灰色
	ColorOrange  = [3]float64{1, 0.5, 0}     // 橙色
	ColorPurple  = [3]float64{0.5, 0, 0.5}   // 紫色
	ColorBrown   = [3]float64{0.6, 0.3, 0}   // 棕色
)

// Color 表示RGBA颜色
type Color struct {
	R, G, B float64 // RGB值，范围0-1
	A       float64 // 透明度，范围0-1，0完全透明，1不透明
}

// NewColor 创建一个新的颜色
func NewColor(r, g, b, a float64) Color {
	return Color{
		R: clamp(r, 0, 1),
		G: clamp(g, 0, 1),
		B: clamp(b, 0, 1),
		A: clamp(a, 0, 1),
	}
}

// ColorToRGBA 将传统的[3]float64颜色转换为Color
func ColorToRGBA(color [3]float64, alpha float64) Color {
	return Color{
		R: color[0],
		G: color[1],
		B: color[2],
		A: alpha,
	}
}

// clamp 将值限制在指定范围内
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// BaseShape 包含所有图形的基本属性
type BaseShape struct {
	ShapeType string     `json:"type"`      // 图形类型
	Points    []*Point   `json:"points"`    // 点集合
	Color     [3]float64 `json:"color"`     // 图形颜色，RGB值(0-1)
	LineWidth float64    `json:"linewidth"` // 线宽
}

// GetType 返回图形的类型
func (b BaseShape) GetType() string {
	return b.ShapeType
}

// GetColor 返回图形的颜色
func (b BaseShape) GetColor() [3]float64 {
	return b.Color
}

// GetPoints 返回图形的点集合
func (b BaseShape) GetPoints() []*Point {
	return b.Points
}

// Option 是图形设置的函数选项接口
type Option func(interface{})

// WithColor 设置图形颜色
func WithColor(color [3]float64) Option {
	return func(s interface{}) {
		if shape, ok := s.(interface{ SetColor([3]float64) }); ok {
			shape.SetColor(color)
		}
	}
}

// WithLineWidth 设置线宽
func WithLineWidth(width float64) Option {
	return func(s interface{}) {
		if shape, ok := s.(interface{ SetLineWidth(float64) }); ok {
			shape.SetLineWidth(width)
		}
	}
}

// WithPoints 设置图形的点集合
func WithPoints(points []*Point) Option {
	return func(s interface{}) {
		if shape, ok := s.(interface{ SetPoints([]*Point) }); ok {
			shape.SetPoints(points)
		}
	}
}

// ShapeFactory 是创建各种图形的工厂接口
type ShapeFactory interface {
	Create(options ...Option) Shape
}
