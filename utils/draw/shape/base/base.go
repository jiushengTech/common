package base

import (
	"image/color"

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
	GetColor() *color.RGBA

	// GetPoints 返回图形的点集合
	GetPoints() []*Point
}

// 颜色常量 (RGB值，范围0-1)
var (
	ColorWhite   = &color.RGBA{R: 255, G: 255, B: 255, A: 255} // 不透明白色
	ColorBlack   = &color.RGBA{R: 0, G: 0, B: 0, A: 255}       // 不透明黑色
	ColorRed     = &color.RGBA{R: 255, G: 0, B: 0, A: 255}     // 不透明红色
	ColorBlue    = &color.RGBA{R: 0, G: 0, B: 255, A: 255}     // 不透明蓝色
	ColorGreen   = &color.RGBA{R: 0, G: 255, B: 0, A: 255}     // 不透明绿色
	ColorYellow  = &color.RGBA{R: 255, G: 255, B: 0, A: 255}   // 黄色
	ColorCyan    = &color.RGBA{R: 0, G: 255, B: 255, A: 255}   // 不透明青色
	ColorMagenta = &color.RGBA{R: 255, G: 0, B: 255, A: 255}   // 不透明品红
	ColorGray    = &color.RGBA{R: 128, G: 128, B: 128, A: 255} // 灰色
	ColorOrange  = &color.RGBA{R: 255, G: 165, B: 0, A: 255}   // 橙色
	ColorPurple  = &color.RGBA{R: 128, G: 0, B: 128, A: 255}   // 紫色
	ColorBrown   = &color.RGBA{R: 139, G: 69, B: 19, A: 255}   // 棕色

	ColorGrayTranslucent  = &color.RGBA{R: 128, G: 128, B: 128, A: 127} // 灰色（半透明）
	ColorBlueTranslucent  = &color.RGBA{R: 0, G: 0, B: 255, A: 127}     // 蓝色（半透明）
	ColorRedTranslucent   = &color.RGBA{R: 255, G: 0, B: 0, A: 127}     // 红色（半透明）
	ColorGreenTranslucent = &color.RGBA{R: 0, G: 255, B: 0, A: 127}     // 绿色（半透明）
	ColorBlackTranslucent = &color.RGBA{R: 0, G: 0, B: 0, A: 127}       // 黑色（半透明）
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
	ShapeType string      `json:"type"`      // 图形类型
	Points    []*Point    `json:"points"`    // 点集合
	Color     *color.RGBA `json:"color"`     // 图形颜色，RGB值(0-1)
	LineWidth float64     `json:"linewidth"` // 线宽
}

// GetType 返回图形的类型
func (b BaseShape) GetType() string {
	return b.ShapeType
}

// GetColor 返回图形的颜色
func (b BaseShape) GetColor() *color.RGBA {
	return b.Color
}

// GetPoints 返回图形的点集合
func (b BaseShape) GetPoints() []*Point {
	return b.Points
}

// Option 是图形设置的函数选项接口
type Option func(interface{})

// WithColor 设置图形颜色
func WithColor(c *color.RGBA) Option {
	return func(s interface{}) {
		if shape, ok := s.(interface{ SetColor(rgba *color.RGBA) }); ok {
			shape.SetColor(c)
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
