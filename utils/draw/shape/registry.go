package shape

import (
	"github.com/jiushengTech/common/utils/draw/shape/base"
	"github.com/jiushengTech/common/utils/draw/shape/group"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/circle"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/hollowpolygon"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/line"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/polygon"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/rectangle"
)

// ShapeRegistry 图形注册表，用于注册和获取图形工厂
type ShapeRegistry struct {
	factories map[string]base.ShapeFactory
}

// NewRegistry 创建新的图形注册表
func NewRegistry() *ShapeRegistry {
	return &ShapeRegistry{
		factories: make(map[string]base.ShapeFactory),
	}
}

// Register 注册图形工厂
func (r *ShapeRegistry) Register(shapeType string, factory base.ShapeFactory) {
	r.factories[shapeType] = factory
}

// GetFactory 获取指定类型的图形工厂
func (r *ShapeRegistry) GetFactory(shapeType string) (base.ShapeFactory, bool) {
	factory, exists := r.factories[shapeType]
	return factory, exists
}

// DefaultRegistry 创建默认的图形注册表并注册所有内置图形
func DefaultRegistry() *ShapeRegistry {
	registry := NewRegistry()

	// 注册垂直线工厂
	registry.Register("vertical_line", line.Factory{LineType: line.Vertical})

	// 注册水平线工厂
	registry.Register("horizontal_line", line.Factory{LineType: line.Horizontal})

	// 注册矩形工厂
	registry.Register("rectangle", rectangle.Factory{})

	// 注册圆形工厂
	registry.Register("circle", circle.Factory{})

	// 注册多边形工厂
	registry.Register("polygon", polygon.Factory{})

	// 注册镂空多边形工厂
	registry.Register("hollow_polygon", hollowpolygon.Factory{})

	// 注册图形组工厂
	registry.Register("shape_group", group.Factory{Name: "default_group"})

	return registry
}

// CreateShape 使用注册表创建指定类型的图形
func CreateShape(registry *ShapeRegistry, shapeType string, options ...base.Option) (base.Shape, bool) {
	factory, exists := registry.GetFactory(shapeType)
	if !exists {
		return nil, false
	}

	return factory.Create(options...), true
}
