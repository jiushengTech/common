package linedraw

import (
	"sort"
)

// NewLine 创建一条线
func NewLine(lineType LineType, points []Point, values []float64, options ...LineOption) Line {
	line := Line{
		Type:      lineType,
		Points:    points,
		Values:    values,
		Color:     ColorYellow, // 默认黄色
		LineWidth: 2.0,
	}

	// 应用所有选项
	for _, option := range options {
		option(&line)
	}

	return line
}

// NewVerticalLine 创建一条竖线
func NewVerticalLine(xpoints []Point, values []float64, options ...LineOption) Line {
	// 按X坐标排序
	sort.Slice(xpoints, func(i, j int) bool {
		return xpoints[i].X < xpoints[j].X
	})
	return NewLine(VerticalLine, xpoints, values, options...)
}

// NewHorizontalLine 创建一条横线
func NewHorizontalLine(ypoints []Point, values []float64, options ...LineOption) Line {
	// 按Y坐标排序
	sort.Slice(ypoints, func(i, j int) bool {
		return ypoints[i].Y < ypoints[j].Y
	})
	return NewLine(HorizontalLine, ypoints, values, options...)
}
