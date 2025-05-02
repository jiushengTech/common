package path

import (
	"github.com/fogleman/gg"
	"github.com/jiushengTech/common/draw/shape/base"
)

// Path 表示图形路径，可以包含多个图形
type Path struct {
	base.BaseShape
	Name    string       `json:"name"`    // 路径名称
	Shapes  []base.Shape `json:"shapes"`  // 包含的图形
	Visible bool         `json:"visible"` // 是否可见
}

// SetColor 设置颜色 - 对Path无效，但需要实现接口
func (p *Path) SetColor(color [3]float64) {
	p.Color = color
}

// SetLineWidth 设置线宽 - 对Path无效，但需要实现接口
func (p *Path) SetLineWidth(width float64) {
	p.LineWidth = width
}

// SetPoints 设置点集合 - 对Path无效，但需要实现接口
func (p *Path) SetPoints(points []base.Point) {
	p.Points = points
}

// New 创建新的图形路径
func New(name string, options ...base.Option) *Path {
	path := &Path{
		BaseShape: base.BaseShape{
			ShapeType: "path",
			Color:     base.ColorWhite,
			LineWidth: 1.0,
		},
		Name:    name,
		Shapes:  []base.Shape{},
		Visible: true,
	}

	// 应用所有选项
	for _, option := range options {
		option(path)
	}

	return path
}

// AddShape 添加图形到路径
func (p *Path) AddShape(shape base.Shape) *Path {
	p.Shapes = append(p.Shapes, shape)
	return p
}

// AddShapes 添加多个图形到路径
func (p *Path) AddShapes(shapes []base.Shape) *Path {
	p.Shapes = append(p.Shapes, shapes...)
	return p
}

// SetVisible 设置路径是否可见
func (p *Path) SetVisible(visible bool) *Path {
	p.Visible = visible
	return p
}

// Draw 绘制路径中所有可见的图形
func (p *Path) Draw(dc *gg.Context, width, height float64) error {
	if !p.Visible {
		return nil
	}

	for _, shape := range p.Shapes {
		if err := shape.Draw(dc, width, height); err != nil {
			return err
		}
	}

	return nil
}

// Factory 创建路径的工厂
type Factory struct {
	Name string // 路径名称
}

// Create 创建路径
func (f Factory) Create(options ...base.Option) base.Shape {
	path := &Path{
		BaseShape: base.BaseShape{
			ShapeType: "path",
			Color:     base.ColorWhite,
			LineWidth: 1.0,
		},
		Name:    f.Name,
		Shapes:  []base.Shape{},
		Visible: true,
	}

	// 应用所有选项
	for _, option := range options {
		option(path)
	}

	return path
}

// WithVisible 设置路径是否可见
func WithVisible(visible bool) base.Option {
	return func(s interface{}) {
		if path, ok := s.(*Path); ok {
			path.Visible = visible
		}
	}
}

// WithName 设置路径名称
func WithName(name string) base.Option {
	return func(s interface{}) {
		if path, ok := s.(*Path); ok {
			path.Name = name
		}
	}
}
