package draw

import (
	"time"

	"github.com/jiushengTech/common/utils/draw/processor"
	"github.com/jiushengTech/common/utils/draw/shape"
	"github.com/jiushengTech/common/utils/draw/shape/base"
	"github.com/jiushengTech/common/utils/draw/shape/group"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/circle"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/line"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/rectangle"
)

// 导出公共类型别名
type (
	// Point 表示二维坐标点
	Point = base.Point

	// ImageProcessor 图像处理器
	ImageProcessor = processor.ImageProcessor

	// ShapeRegistry 图形注册表类型
	ShapeRegistry = shape.ShapeRegistry

	// LineType 线条类型
	LineType = line.Type

	// Color 表示RGBA颜色
	Color = base.Color
)

// 线条类型常量
const (
	VerticalLine   = line.Vertical   // 竖线
	HorizontalLine = line.Horizontal // 横线
)

// 图像输出格式
const (
	FormatPNG  = processor.FormatPNG  // PNG格式
	FormatJPEG = processor.FormatJPEG // JPEG格式
)

// 颜色常量
var (
	ColorWhite   = base.ColorWhite   // 白色
	ColorBlack   = base.ColorBlack   // 黑色
	ColorRed     = base.ColorRed     // 红色
	ColorBlue    = base.ColorBlue    // 蓝色
	ColorGreen   = base.ColorGreen   // 绿色
	ColorYellow  = base.ColorYellow  // 黄色
	ColorCyan    = base.ColorCyan    // 青色
	ColorMagenta = base.ColorMagenta // 品红
	ColorGray    = base.ColorGray    // 灰色
	ColorOrange  = base.ColorOrange  // 橙色
	ColorPurple  = base.ColorPurple  // 紫色
	ColorBrown   = base.ColorBrown   // 棕色
)

// Registry 全局图形注册表实例
var Registry = shape.DefaultRegistry()

// 图形创建函数
// -------------------------

// NewShape 通过类型名称创建图形
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
func NewRectangle(topLeft, bottomRight Point, options ...ShapeOption) Shape {
	points := []Point{topLeft, bottomRight}
	return rectangle.New(points, options...)
}

// NewCircle 创建一个圆形
func NewCircle(center Point, radius float64, options ...ShapeOption) Shape {
	return circle.New(center, radius, options...)
}

// NewShapeGroup 创建一个图形组
func NewShapeGroup(name string, options ...ShapeOption) *ShapeGroup {
	return group.New(name, options...)
}

// NewShapeRegistry 创建一个新的图形注册表
func NewShapeRegistry() *ShapeRegistry {
	return shape.NewRegistry()
}

// RegisterShape 注册自定义图形到全局注册表
func RegisterShape(shapeType string, factory base.ShapeFactory) {
	Registry.Register(shapeType, factory)
}

// NewImageProcessor 创建一个新的图像处理器
func NewImageProcessor(imagePath string, options ...ProcessorOption) *ImageProcessor {
	return processor.NewImageProcessor(imagePath, options...)
}

// 颜色相关函数
// -------------------------

// NewColor 创建一个新的RGBA颜色
func NewColor(r, g, b, a float64) Color {
	return base.NewColor(r, g, b, a)
}

// ColorToRGBA 将RGB颜色转换为RGBA颜色
func ColorToRGBA(color [3]float64, alpha float64) Color {
	return base.ColorToRGBA(color, alpha)
}

// 图形配置选项函数
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

// WithVisible 设置图形组是否可见
func WithVisible(visible bool) ShapeOption {
	return group.WithVisible(visible)
}

// WithName 设置图形组名称
func WithName(name string) ShapeOption {
	return group.WithName(name)
}

// 处理器配置选项函数
// -------------------------

// WithOutputName 设置输出文件名
func WithOutputName(name string) ProcessorOption {
	return processor.WithOutputName(name)
}

// WithCustomName 设置自定义名称（兼容）
func WithCustomName(name string) ProcessorOption {
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
func WithOutputFormat(format OutputFormat) ProcessorOption {
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

// WithPreProcess 设置图像预处理函数
func WithPreProcess(fn processor.ProcessFunc) ProcessorOption {
	return processor.WithPreProcess(fn)
}

// WithPostProcess 设置图像后处理函数
func WithPostProcess(fn processor.ProcessFunc) ProcessorOption {
	return processor.WithPostProcess(fn)
}

// CleanupAllTempFiles 清理所有临时文件
func CleanupAllTempFiles() {
	processor.CleanupAllTempFiles()
}

// GetAbsoluteOutputPath 获取处理器的绝对输出路径
func GetAbsoluteOutputPath(p *ImageProcessor) (string, error) {
	return p.GetAbsoluteOutputPath()
}

// ProcessImage 处理图像并返回结果文件路径
func ProcessImage(p *ImageProcessor) (string, error) {
	return p.Process()
}
