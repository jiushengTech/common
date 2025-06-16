package line

import (
	"image/color"

	"github.com/jiushengTech/common/utils/drawutil/shape/base"
)

// LineOptions 线条的配置选项
type LineOptions struct {
	Type         Type      `json:"line_type"`     // 线条类型
	Values       []float64 `json:"values"`        // 点之间的值（长度比点少1）
	TextPosition float64   `json:"text_position"` // 文本位置(0-1之间的值，表示在两条线之间的位置比例)
}

// LineOption 线条配置选项函数类型
type LineOption func(*LineOptions, *base.BaseShape)

// WithType 设置线条类型
func WithType(lineType Type) LineOption {
	return func(opts *LineOptions, shape *base.BaseShape) {
		opts.Type = lineType
	}
}

// WithValues 设置值集合
func WithValues(values []float64) LineOption {
	return func(opts *LineOptions, shape *base.BaseShape) {
		opts.Values = values
	}
}

// WithTextPosition 设置文本位置
func WithTextPosition(position float64) LineOption {
	return func(opts *LineOptions, shape *base.BaseShape) {
		opts.TextPosition = position
	}
}

// WithPoints 设置点集合
func WithPoints(points []*base.Point) LineOption {
	return func(opts *LineOptions, shape *base.BaseShape) {
		shape.Points = points
	}
}

// WithColor 设置颜色
func WithColor(c *color.RGBA) LineOption {
	return func(opts *LineOptions, shape *base.BaseShape) {
		shape.Color = c
	}
}

// WithLineWidth 设置线宽
func WithLineWidth(width float64) LineOption {
	return func(opts *LineOptions, shape *base.BaseShape) {
		shape.LineWidth = width
	}
}
