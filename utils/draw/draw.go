// Package draw 提供图像绘制功能，支持多种形状和图像处理操作
package draw

import (
	"time"

	"github.com/jiushengTech/common/utils/draw/processor"
	"github.com/jiushengTech/common/utils/draw/shape"
	"github.com/jiushengTech/common/utils/draw/shape/base"
	"github.com/jiushengTech/common/utils/draw/shape/path"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/circle"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/line"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/rectangle"
)

// 导出公共类型
type (
	// Point 表示二维坐标点
	Point = base.Point

	// Shape 是所有图形的通用接口
	Shape = base.Shape

	// ImageProcessor 图像处理器
	ImageProcessor = processor.ImageProcessor

	// ShapeRegistry 图形注册表类型
	ShapeRegistry = shape.ShapeRegistry

	// LineType 线条类型
	LineType = line.Type

	// Path 路径类型
	Path = path.Path

	// ShapeOption 图形配置选项类型
	ShapeOption = base.Option

	// ProcessorOption 处理器配置选项类型
	ProcessorOption = processor.Option
)

// 线条类型常量
const (
	VerticalLine   = line.Vertical   // 竖线
	HorizontalLine = line.Horizontal // 横线
)

// 颜色常量
var (
	ColorWhite  = base.ColorWhite  // 白色
	ColorBlack  = base.ColorBlack  // 黑色
	ColorRed    = base.ColorRed    // 红色
	ColorBlue   = base.ColorBlue   // 蓝色
	ColorGreen  = base.ColorGreen  // 绿色
	ColorYellow = base.ColorYellow // 黄色

	// 可以添加更多预定义颜色...
)

// 默认输出文件名
const DefaultOutputName = processor.DefaultOutputName

// 获取默认输出文件名（基于当前时间格式）
func GetDefaultOutputName(format processor.OutputFormat) string {
	return processor.GetDefaultOutputName(format)
}

// Registry 全局图形注册表实例
var Registry = shape.DefaultRegistry()

// 图形创建函数
// -------------------------

// NewShape 通过类型名称创建图形
// 示例:
//
//	circle, ok := draw.NewShape("circle",
//	   draw.WithPoints([]draw.Point{{X: 400, Y: 400}}),
//	   draw.WithRadius(60),
//	   draw.WithColor(draw.ColorRed),
//	)
func NewShape(shapeType string, options ...ShapeOption) (Shape, bool) {
	return shape.CreateShape(Registry, shapeType, options...)
}

// NewLine 创建一条线
func NewLine(lineType LineType, points []Point, values []float64, options ...ShapeOption) Shape {
	return line.New(lineType, points, values, options...)
}

// NewVerticalLine 创建一条竖线
func NewVerticalLine(xpoints []Point, values []float64, options ...ShapeOption) Shape {
	factory := line.Factory{LineType: VerticalLine}
	options = append(options, base.WithPoints(xpoints), line.WithValues(values))
	return factory.Create(options...)
}

// NewHorizontalLine 创建一条横线
func NewHorizontalLine(ypoints []Point, values []float64, options ...ShapeOption) Shape {
	factory := line.Factory{LineType: HorizontalLine}
	options = append(options, base.WithPoints(ypoints), line.WithValues(values))
	return factory.Create(options...)
}

// NewRectangle 创建一个矩形
// 参数:
//
//	topLeft: 左上角坐标
//	bottomRight: 右下角坐标
//	options: 可选配置，如颜色、线宽、是否填充等
func NewRectangle(topLeft, bottomRight Point, options ...ShapeOption) Shape {
	points := []Point{topLeft, bottomRight}
	return rectangle.New(points, options...)
}

// NewCircle 创建一个圆形
// 参数:
//
//	center: 圆心坐标
//	radius: 半径
//	options: 可选配置，如颜色、线宽、是否填充等
func NewCircle(center Point, radius float64, options ...ShapeOption) Shape {
	return circle.New(center, radius, options...)
}

// NewPath 创建一个图形路径
func NewPath(name string, options ...ShapeOption) *Path {
	return path.New(name, options...)
}

// NewShapeRegistry 创建一个新的图形注册表
// 可用于创建自定义的图形类型注册表
func NewShapeRegistry() *ShapeRegistry {
	return shape.NewRegistry()
}

// RegisterShape 注册自定义图形到全局注册表
// 参数:
//
//	shapeType: 图形类型名称
//	factory: 图形工厂实例
func RegisterShape(shapeType string, factory base.ShapeFactory) {
	Registry.Register(shapeType, factory)
}

// NewImageProcessor 创建一个新的图像处理器
// 参数:
//
//	imagePath: 图像路径，可以是本地文件路径或URL
//	options: 可选配置，如输出目录、文件名、添加图形等
func NewImageProcessor(imagePath string, options ...ProcessorOption) *ImageProcessor {
	return processor.NewImageProcessor(imagePath, options...)
}

// 图形配置选项
// -------------------------

// WithColor 设置图形颜色
func WithColor(color [3]float64) ShapeOption {
	return base.WithColor(color)
}

// WithLineWidth 设置线宽
func WithLineWidth(width float64) ShapeOption {
	return base.WithLineWidth(width)
}

// WithPoints 设置图形的点集合
func WithPoints(points []Point) ShapeOption {
	return base.WithPoints(points)
}

// WithValues 设置线条的值集合
func WithValues(values []float64) ShapeOption {
	return line.WithValues(values)
}

// WithLineType 设置线条类型
func WithLineType(lineType LineType) ShapeOption {
	return line.WithLineType(lineType)
}

// WithTextPosition 设置线条文本位置 (0-1之间)
// 0表示最左/最上方，0.5表示中间，1表示最右/最下方
func WithTextPosition(position float64) ShapeOption {
	return line.WithTextPosition(position)
}

// WithFill 设置是否填充
func WithFill(fill bool) ShapeOption {
	return rectangle.WithFill(fill)
}

// WithRadius 设置圆的半径
func WithRadius(radius float64) ShapeOption {
	return circle.WithRadius(radius)
}

// WithVisible 设置路径是否可见
func WithVisible(visible bool) ShapeOption {
	return path.WithVisible(visible)
}

// WithName 设置路径名称
func WithName(name string) ShapeOption {
	return path.WithName(name)
}

// 处理器配置选项
// -------------------------

// WithOutputName 设置输出文件名
func WithOutputName(name string) ProcessorOption {
	return processor.WithOutputName(name)
}

// WithTimeBasedName 设置输出文件名为当前时间
func WithTimeBasedName() ProcessorOption {
	return processor.WithTimeBasedName()
}

// WithOutputDir 设置输出目录
func WithOutputDir(dir string) ProcessorOption {
	return processor.WithOutputDir(dir)
}

// WithShape 添加一个图形
func WithShape(shape Shape) ProcessorOption {
	return processor.WithShape(shape)
}

// WithShapes 添加多个图形
func WithShapes(shapes []Shape) ProcessorOption {
	return processor.WithShapes(shapes)
}

// WithOutputFormat 设置输出格式
func WithOutputFormat(format processor.OutputFormat) ProcessorOption {
	return processor.WithOutputFormat(format)
}

// WithJpegQuality 设置JPEG质量
func WithJpegQuality(quality int) ProcessorOption {
	return processor.WithJpegQuality(quality)
}

// WithRequestTimeout 设置HTTP请求超时时间
func WithRequestTimeout(timeout time.Duration) ProcessorOption {
	return processor.WithRequestTimeout(timeout)
}

// WithPreProcess 设置预处理函数
func WithPreProcess(fn processor.ProcessFunc) ProcessorOption {
	return processor.WithPreProcess(fn)
}

// WithPostProcess 设置后处理函数
func WithPostProcess(fn processor.ProcessFunc) ProcessorOption {
	return processor.WithPostProcess(fn)
}

// CleanupAllTempFiles 清理所有临时文件和目录
func CleanupAllTempFiles() {
	processor.CleanupAllTempFiles()
}

// GetAbsoluteOutputPath 获取图片输出的绝对路径
func GetAbsoluteOutputPath(p *ImageProcessor) (string, error) {
	return p.GetAbsoluteOutputPath()
}

// ProcessImage 处理图像，返回生成图片的绝对路径
func ProcessImage(p *ImageProcessor) (string, error) {
	return p.Process()
}

// 导出输出格式
var (
	FormatPNG  = processor.FormatPNG
	FormatJPEG = processor.FormatJPEG
)
