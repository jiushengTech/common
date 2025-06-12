package polygonops

import (
	"image/color"

	"github.com/jiushengTech/common/utils/drawutil/shape/base"
)

// Option 定义多边形操作的选项函数类型
type Option func(*PolygonOperation)

// WithFillColor 设置填充颜色
func WithFillColor(color *color.RGBA) Option {
	return func(po *PolygonOperation) {
		po.FillColor = color
	}
}

// WithOutlineColor 设置轮廓颜色
func WithOutlineColor(color *color.RGBA) Option {
	return func(po *PolygonOperation) {
		po.OutlineColor = color
	}
}

// WithOutlineWidth 设置轮廓宽度
func WithOutlineWidth(width float64) Option {
	return func(po *PolygonOperation) {
		po.OutlineWidth = width
	}
}

// WithDrawOutline 设置是否绘制轮廓
func WithDrawOutline(draw bool) Option {
	return func(po *PolygonOperation) {
		po.DrawOutline = draw
	}
}

// WithPolygonA 设置多边形A的点
func WithPolygonA(points []*base.Point) Option {
	return func(po *PolygonOperation) {
		po.PolygonA = points
	}
}

// WithPolygonB 设置多边形B的点
func WithPolygonB(points []*base.Point) Option {
	return func(po *PolygonOperation) {
		po.PolygonB = points
	}
}

// WithOperationType 设置操作类型
func WithOperationType(opType OperationType) Option {
	return func(po *PolygonOperation) {
		po.Type = opType
	}
}

// DrawOption 定义绘制多边形操作的选项函数类型
type DrawOption func(*DrawConfig)

// DrawConfig 定义绘制多边形操作的配置
type DrawConfig struct {
	// 图像URL
	ImageURL string
	// 多边形操作列表
	Operations []*PolygonOperation
	// 输出目录
	OutputDir string
	// 输出文件名
	OutputName string
	// 其他可扩展选项
	// ...
}

// WithImageURL 设置图像URL
func WithImageURL(url string) DrawOption {
	return func(cfg *DrawConfig) {
		cfg.ImageURL = url
	}
}

// WithOperations 设置多边形操作列表
func WithOperations(operations []*PolygonOperation) DrawOption {
	return func(cfg *DrawConfig) {
		cfg.Operations = operations
	}
}

// WithOutputDir 设置输出目录
func WithOutputDirectory(dir string) DrawOption {
	return func(cfg *DrawConfig) {
		cfg.OutputDir = dir
	}
}

// WithOutputFileName 设置输出文件名
func WithOutputFileName(name string) DrawOption {
	return func(cfg *DrawConfig) {
		cfg.OutputName = name
	}
}

// AddOperation 添加单个多边形操作
func AddOperation(operation *PolygonOperation) DrawOption {
	return func(cfg *DrawConfig) {
		cfg.Operations = append(cfg.Operations, operation)
	}
}
