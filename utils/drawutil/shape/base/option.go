package base

import "image/color"

// Option 是 BaseShape 配置选项
type Option func(*BaseShape)

// WithColor 设置图形颜色
func WithColor(c *color.RGBA) Option {
	return func(s *BaseShape) {
		s.Color = c
	}
}

// WithLineWidth 设置线宽
func WithLineWidth(w float64) Option {
	return func(s *BaseShape) {
		s.LineWidth = w
	}
}

// WithPoints 设置点集合
func WithPoints(p []*Point) Option {
	return func(s *BaseShape) {
		s.Points = p
	}
}

// WithShapeType 设置图形类型
func WithShapeType(t string) Option {
	return func(s *BaseShape) {
		s.ShapeType = t
	}
}

// ApplyOptions 应用选项到 BaseShape
func ApplyOptions(base *BaseShape, options ...Option) {
	for _, option := range options {
		option(base)
	}
}
