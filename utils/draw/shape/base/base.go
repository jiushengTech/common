package base

import (
	"image/color"

	"github.com/fogleman/gg"
)

//// ---------- 基础类型 ----------

// Point 表示二维坐标点
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Shape 是所有图形的通用接口
type Shape interface {
	Draw(dc *gg.Context, width, height float64) error
	GetType() string
	GetColor() *color.RGBA
	GetPoints() []*Point
}

// BaseShape 包含所有图形的基本属性
type BaseShape struct {
	ShapeType string      `json:"type"`      // 图形类型
	Points    []*Point    `json:"points"`    // 点集合
	Color     *color.RGBA `json:"color"`     // 图形颜色
	LineWidth float64     `json:"lineWidth"` // 线宽
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
