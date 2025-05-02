// Package draw 提供图像绘制功能，支持多种形状和图像处理操作
package draw

import (
	"time"

	"github.com/jiushengTech/common/utils/draw/processor"
	"github.com/jiushengTech/common/utils/draw/shape"
	"github.com/jiushengTech/common/utils/draw/shape/base"
	"github.com/jiushengTech/common/utils/draw/shape/group"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/circle"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/hollowpolygon"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/line"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/polygon"
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

	// ShapeGroup 图形组类型
	ShapeGroup = group.ShapeGroup

	// ShapeOption 图形配置选项类型
	ShapeOption = base.Option

	// ProcessorOption 处理器配置选项类型
	ProcessorOption = processor.Option

	// OutputFormat 输出格式类型
	OutputFormat = processor.OutputFormat

	// Color 表示RGBA颜色
	Color = base.Color

	// ProcessFunc 图像处理函数类型
	ProcessFunc = processor.ProcessFunc
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

// 默认输出文件名
const DefaultOutputName = processor.DefaultOutputName

// GetDefaultOutputName 获取默认输出文件名（基于当前时间格式）
func GetDefaultOutputName(format OutputFormat) string {
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
	// 直接调用底层包的构造函数
	return line.New(lineType, points, values, options...)
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

// NewShapeGroup 创建一个图形组
func NewShapeGroup(name string, options ...ShapeOption) *ShapeGroup {
	return group.New(name, options...)
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

// NewPolygon 创建一个多边形
// 参数:
//
//	points: 多边形的顶点坐标数组，至少需要3个点
//	options: 可选配置，如颜色、线宽、是否填充等
func NewPolygon(points []Point, options ...ShapeOption) Shape {
	return polygon.New(points, options...)
}

// NewHollowPolygon 创建一个镂空多边形
// 参数:
//
//	outerPoints: 外部多边形的顶点数组，至少需要3个点
//	innerPoints: 内部多边形的顶点数组，至少需要3个点
//	options: 可选配置，如颜色、线宽、不透明度等
func NewHollowPolygon(outerPoints, innerPoints []Point, options ...ShapeOption) Shape {
	return hollowpolygon.New(outerPoints, innerPoints, options...)
}

// 颜色相关函数
// -------------------------

// NewColor 创建一个新的RGBA颜色
// 参数:
//
//	r,g,b: RGB颜色值(0-1)
//	a: 透明度(0-1)，0完全透明，1不透明
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

// WithVisible 设置图形组是否可见
func WithVisible(visible bool) ShapeOption {
	return group.WithVisible(visible)
}

// WithName 设置图形组名称
func WithName(name string) ShapeOption {
	return group.WithName(name)
}

// WithPolygonFill 设置是否填充多边形
func WithPolygonFill(fill bool) ShapeOption {
	return polygon.WithFill(fill)
}

// WithHollowPolygonOpacity 设置镂空多边形的不透明度
func WithHollowPolygonOpacity(opacity float64) ShapeOption {
	return hollowpolygon.WithOpacity(opacity)
}

// WithOuterPoints 设置镂空多边形的外部顶点
func WithOuterPoints(points []Point) ShapeOption {
	return hollowpolygon.WithOuterPoints(points)
}

// WithInnerPoints 设置镂空多边形的内部顶点
func WithInnerPoints(points []Point) ShapeOption {
	return hollowpolygon.WithInnerPoints(points)
}

// 处理器配置选项函数
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
func WithPreProcess(fn ProcessFunc) ProcessorOption {
	return processor.WithPreProcess(fn)
}

// WithPostProcess 设置图像后处理函数
func WithPostProcess(fn ProcessFunc) ProcessorOption {
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
