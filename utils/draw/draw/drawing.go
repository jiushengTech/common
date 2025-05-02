package draw

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jiushengTech/common/utils/draw/processor"
	"github.com/jiushengTech/common/utils/draw/shape/base"
	"github.com/jiushengTech/common/utils/draw/shape/group"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/circle"
	"github.com/jiushengTech/common/utils/draw/shape/primitives/rectangle"
)

// 定义类型别名
type (
	// Shape 是所有图形的通用接口
	Shape = base.Shape

	// ShapeGroup 图形组类型
	ShapeGroup = group.ShapeGroup

	// ShapeOption 图形配置选项类型
	ShapeOption = base.Option

	// ProcessorOption 处理器配置选项类型
	ProcessorOption = processor.Option

	// OutputFormat 输出格式类型
	OutputFormat = processor.OutputFormat
)

// Drawing 表示一个绘图会话，使用流畅的API设计
type Drawing struct {
	// 基本配置
	imageURL         string
	outputDir        string
	filename         string
	useTimeBasedName bool
	format           OutputFormat
	quality          int // JPEG压缩质量

	// 内容配置
	shapes         []Shape
	shapeGroups    map[string]*ShapeGroup
	preProcessors  []processor.ProcessFunc // 使用 processor 包的类型
	postProcessors []processor.ProcessFunc // 使用 processor 包的类型
}

// ShapeOptions 包含形状的基本选项
type ShapeOptions struct {
	Color     [3]float64
	LineWidth float64
	Fill      bool
}

// NewDrawing 创建一个新的绘图会话
func NewDrawing(imageURL string) *Drawing {
	return &Drawing{
		imageURL:       imageURL,
		outputDir:      "output", // 默认输出目录
		filename:       "image",  // 默认文件名
		shapes:         make([]Shape, 0),
		shapeGroups:    make(map[string]*ShapeGroup),
		preProcessors:  make([]processor.ProcessFunc, 0),
		postProcessors: make([]processor.ProcessFunc, 0),
		format:         FormatPNG,
		quality:        85, // 默认JPEG质量
	}
}

// WithOutputDirectory 设置输出目录
func (d *Drawing) WithOutputDirectory(dir string) *Drawing {
	d.outputDir = dir
	return d
}

// WithFilename 设置输出文件名
func (d *Drawing) WithFilename(name string) *Drawing {
	d.filename = name
	return d
}

// WithTimeBasedFilename 使用基于时间的文件名
func (d *Drawing) WithTimeBasedFilename() *Drawing {
	d.useTimeBasedName = true
	return d
}

// WithOutputFormat 设置输出格式
func (d *Drawing) WithOutputFormat(format OutputFormat) *Drawing {
	d.format = format
	return d
}

// WithJPEGQuality 设置JPEG压缩质量
func (d *Drawing) WithJPEGQuality(quality int) *Drawing {
	if quality < 1 {
		quality = 1
	}
	if quality > 100 {
		quality = 100
	}
	d.quality = quality
	return d
}

// AddShape 添加一个形状
func (d *Drawing) AddShape(shape Shape) *Drawing {
	d.shapes = append(d.shapes, shape)
	return d
}

// AddCircle 添加一个圆形
func (d *Drawing) AddCircle(center Point, radius float64, opts ...ShapeOption) *Drawing {
	circle := NewCircle(center, radius, opts...)
	return d.AddShape(circle)
}

// AddRectangle 添加一个矩形
func (d *Drawing) AddRectangle(topLeft, bottomRight Point, opts ...ShapeOption) *Drawing {
	rect := NewRectangle(topLeft, bottomRight, opts...)
	return d.AddShape(rect)
}

// AddLine 添加一条线
func (d *Drawing) AddLine(start, end Point, opts ...ShapeOption) *Drawing {
	// 创建一条直线
	points := []Point{start, end}
	lineType := VerticalLine
	if start.X != end.X && start.Y == end.Y {
		lineType = HorizontalLine
	}
	line := NewLine(lineType, points, nil, opts...)
	return d.AddShape(line)
}

// AddShapeGroup 添加一个形状组
func (d *Drawing) AddShapeGroup(name string, configurator func(*ShapeGroup)) *Drawing {
	group := NewShapeGroup(name)
	configurator(group)
	d.shapeGroups[name] = group

	// 也将组中的所有形状添加到形状列表中
	for _, shape := range group.Shapes {
		d.shapes = append(d.shapes, shape)
	}

	return d
}

// WithPreProcess 添加前处理函数
func (d *Drawing) WithPreProcess(fn processor.ProcessFunc) *Drawing {
	d.preProcessors = append(d.preProcessors, fn)
	return d
}

// WithPostProcess 添加后处理函数
func (d *Drawing) WithPostProcess(fn processor.ProcessFunc) *Drawing {
	d.postProcessors = append(d.postProcessors, fn)
	return d
}

// Process 处理图像并保存
func (d *Drawing) Process() (string, error) {
	// 创建输出目录
	if err := os.MkdirAll(d.outputDir, 0755); err != nil {
		return "", fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 生成文件名
	finalFilename := d.filename
	if d.useTimeBasedName {
		now := time.Now()
		finalFilename = now.Format("20060102_150405")
	}

	// 确保有扩展名
	ext := string(d.format)
	if ext == "" {
		ext = "png"
	}

	if filepath.Ext(finalFilename) == "" {
		finalFilename += "." + ext
	}

	// 组装处理器选项
	opts := []ProcessorOption{
		WithOutputName(finalFilename),
		WithOutputDir(d.outputDir),
		WithOutputFormat(d.format),
		WithJpegQuality(d.quality),
	}

	// 添加形状
	for _, shape := range d.shapes {
		opts = append(opts, WithShape(shape))
	}

	// 添加预处理函数
	for _, pre := range d.preProcessors {
		opts = append(opts, WithPreProcess(pre))
	}

	// 添加后处理函数
	for _, post := range d.postProcessors {
		opts = append(opts, WithPostProcess(post))
	}

	// 创建并运行处理器
	processor := NewImageProcessor(d.imageURL, opts...)
	absPath, err := processor.Process()
	if err != nil {
		return "", fmt.Errorf("处理图像失败: %w", err)
	}

	return absPath, nil
}

// CircleBuilder 提供圆形的流畅构建API
type CircleBuilder struct {
	center  Point
	radius  float64
	options ShapeOptions
}

// Circle 创建一个可链式调用的圆形
func Circle(center Point, radius float64) *CircleBuilder {
	return &CircleBuilder{
		center: center,
		radius: radius,
		options: ShapeOptions{
			Color:     ColorBlack,
			LineWidth: 1.0,
			Fill:      false,
		},
	}
}

// WithColor 设置圆的颜色
func (c *CircleBuilder) WithColor(color [3]float64) *CircleBuilder {
	c.options.Color = color
	return c
}

// WithLineWidth 设置圆的线宽
func (c *CircleBuilder) WithLineWidth(width float64) *CircleBuilder {
	c.options.LineWidth = width
	return c
}

// WithFill 设置圆是否填充
func (c *CircleBuilder) WithFill(fill bool) *CircleBuilder {
	c.options.Fill = fill
	return c
}

// Build 构建圆形并转换为Shape接口
func (c *CircleBuilder) Build() Shape {
	return circle.New(c.center, c.radius,
		WithColor(c.options.Color),
		WithLineWidth(c.options.LineWidth),
		WithFill(c.options.Fill),
	)
}

// RectangleBuilder 提供矩形的流畅构建API
type RectangleBuilder struct {
	topLeft, bottomRight Point
	options              ShapeOptions
}

// Rectangle 创建一个可链式调用的矩形
func Rectangle(topLeft Point, width, height float64) *RectangleBuilder {
	return &RectangleBuilder{
		topLeft:     topLeft,
		bottomRight: Point{X: topLeft.X + width, Y: topLeft.Y + height},
		options: ShapeOptions{
			Color:     ColorBlack,
			LineWidth: 1.0,
			Fill:      false,
		},
	}
}

// WithColor 设置矩形的颜色
func (r *RectangleBuilder) WithColor(color [3]float64) *RectangleBuilder {
	r.options.Color = color
	return r
}

// WithLineWidth 设置矩形的线宽
func (r *RectangleBuilder) WithLineWidth(width float64) *RectangleBuilder {
	r.options.LineWidth = width
	return r
}

// WithFill 设置矩形是否填充
func (r *RectangleBuilder) WithFill(fill bool) *RectangleBuilder {
	r.options.Fill = fill
	return r
}

// Build 构建矩形并转换为Shape接口
func (r *RectangleBuilder) Build() Shape {
	points := []Point{r.topLeft, r.bottomRight}
	return rectangle.New(points,
		WithColor(r.options.Color),
		WithLineWidth(r.options.LineWidth),
		WithFill(r.options.Fill),
	)
}
