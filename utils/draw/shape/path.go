package shape

import (
	"github.com/fogleman/gg"
)

// Path 表示图形路径，可以包含多个图形
type Path struct {
	Name    string  `json:"name"`    // 路径名称
	Shapes  []Shape `json:"shapes"`  // 包含的图形
	Visible bool    `json:"visible"` // 是否可见
}

// NewPath 创建新的图形路径
func NewPath(name string) *Path {
	return &Path{
		Name:    name,
		Shapes:  []Shape{},
		Visible: true,
	}
}

// AddShape 添加图形到路径
func (p *Path) AddShape(shape Shape) *Path {
	p.Shapes = append(p.Shapes, shape)
	return p
}

// AddShapes 添加多个图形到路径
func (p *Path) AddShapes(shapes []Shape) *Path {
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

// GetType 实现Shape接口
func (p *Path) GetType() string {
	return "path"
}

// GetColor 实现Shape接口，路径本身没有颜色
func (p *Path) GetColor() [3]float64 {
	return [3]float64{0, 0, 0}
}

// GetPoints 实现Shape接口，路径本身没有点
func (p *Path) GetPoints() []Point {
	return []Point{}
}
