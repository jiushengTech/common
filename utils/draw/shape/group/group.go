package group

import (
	"github.com/fogleman/gg"
	"github.com/jiushengTech/common/utils/draw/shape/base"
)

// ShapeGroup 表示图形组，可以包含多个图形
type ShapeGroup struct {
	base.BaseShape
	Name    string       `json:"name"`    // 图形组名称
	Shapes  []base.Shape `json:"shapes"`  // 包含的图形
	Visible bool         `json:"visible"` // 是否可见
}

// SetColor 设置颜色 - 对ShapeGroup无效，但需要实现接口
func (g *ShapeGroup) SetColor(color [3]float64) {
	g.Color = color
}

// SetLineWidth 设置线宽 - 对ShapeGroup无效，但需要实现接口
func (g *ShapeGroup) SetLineWidth(width float64) {
	g.LineWidth = width
}

// SetPoints 设置点集合 - 对ShapeGroup无效，但需要实现接口
func (g *ShapeGroup) SetPoints(points []base.Point) {
	g.Points = points
}

// New 创建新的图形组
func New(name string, options ...base.Option) *ShapeGroup {
	group := &ShapeGroup{
		BaseShape: base.BaseShape{
			ShapeType: "shape_group",
			Color:     base.ColorWhite,
			LineWidth: 1.0,
		},
		Name:    name,
		Shapes:  []base.Shape{},
		Visible: true,
	}

	// 应用所有选项
	for _, option := range options {
		option(group)
	}

	return group
}

// AddShape 添加图形到图形组
func (g *ShapeGroup) AddShape(shape base.Shape) *ShapeGroup {
	g.Shapes = append(g.Shapes, shape)
	return g
}

// AddShapes 添加多个图形到图形组
func (g *ShapeGroup) AddShapes(shapes []base.Shape) *ShapeGroup {
	g.Shapes = append(g.Shapes, shapes...)
	return g
}

// SetVisible 设置图形组是否可见
func (g *ShapeGroup) SetVisible(visible bool) *ShapeGroup {
	g.Visible = visible
	return g
}

// Draw 绘制图形组中所有可见的图形
func (g *ShapeGroup) Draw(dc *gg.Context, width, height float64) error {
	if !g.Visible {
		return nil
	}

	for _, shape := range g.Shapes {
		if err := shape.Draw(dc, width, height); err != nil {
			return err
		}
	}

	return nil
}

// Factory 创建图形组的工厂
type Factory struct {
	Name string // 图形组名称
}

// Create 创建图形组
func (f Factory) Create(options ...base.Option) base.Shape {
	group := &ShapeGroup{
		BaseShape: base.BaseShape{
			ShapeType: "shape_group",
			Color:     base.ColorWhite,
			LineWidth: 1.0,
		},
		Name:    f.Name,
		Shapes:  []base.Shape{},
		Visible: true,
	}

	// 应用所有选项
	for _, option := range options {
		option(group)
	}

	return group
}

// WithVisible 设置图形组是否可见
func WithVisible(visible bool) base.Option {
	return func(s interface{}) {
		if group, ok := s.(*ShapeGroup); ok {
			group.Visible = visible
		}
	}
}

// WithName 设置图形组名称
func WithName(name string) base.Option {
	return func(s interface{}) {
		if group, ok := s.(*ShapeGroup); ok {
			group.Name = name
		}
	}
}
