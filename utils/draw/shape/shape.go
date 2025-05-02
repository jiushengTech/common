package shape

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
	GetPoints() []Point
}

// 颜色常量
var (
	ColorWhite  = [3]float64{1, 1, 1} // 白色
	ColorBlack  = [3]float64{0, 0, 0} // 黑色
	ColorRed    = [3]float64{1, 0, 0} // 红色
	ColorBlue   = [3]float64{0, 0, 1} // 蓝色
	ColorGreen  = [3]float64{0, 1, 0} // 绿色
	ColorYellow = [3]float64{1, 1, 0} // 黄色
)

// BaseShape 包含所有图形的基本属性
type BaseShape struct {
	ShapeType string     `json:"type"`      // 图形类型
	Points    []Point    `json:"points"`    // 点集合
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
func (b BaseShape) GetPoints() []Point {
	return b.Points
}

// Option 是图形设置的函数选项接口
type Option func(interface{})

// WithColor 设置图形颜色
func WithColor(color [3]float64) Option {
	return func(s interface{}) {
		if shape, ok := s.(*BaseShape); ok {
			shape.Color = color
		}
	}
}

// WithLineWidth 设置线宽
func WithLineWidth(width float64) Option {
	return func(s interface{}) {
		if shape, ok := s.(*BaseShape); ok {
			shape.LineWidth = width
		}
	}
}

// WithPoints 设置图形的点集合
func WithPoints(points []Point) Option {
	return func(s interface{}) {
		if shape, ok := s.(*BaseShape); ok {
			shape.Points = points
		}
	}
}

// ShapeFactory 是创建各种图形的工厂接口
type ShapeFactory interface {
	Create(options ...Option) Shape
}
