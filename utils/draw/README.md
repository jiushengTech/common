# Draw 图形绘制库

这是一个易用的Go语言图形绘制库，支持在图片上绘制各种形状，如线条、矩形、圆形等，支持多种输出格式，并提供灵活的配置选项。

## 功能特点

- 支持多种基本图形：线条（水平线、垂直线）、矩形、圆形等
- 支持自定义图形（通过实现Shape接口）
- 支持多种输出格式（PNG、JPEG）及质量配置
- 支持本地文件和在线URL图片处理
- 支持图像预处理和后处理功能
- 使用函数选项模式，提供灵活易用的API
- 支持图形工厂模式，易于扩展
- 支持临时文件自动清理

## 项目结构

```
draw/
├── draw.go          # 主包，提供公共API
├── line/            # 线条相关实现
├── processor/       # 图像处理器实现
└── example/         # 使用示例
```

## 安装

```bash
go get github.com/jiushengTech/common/draw
```

## 基本用法

### 创建和绘制基本图形

```go
// 创建矩形
rect := draw.NewRectangle(
    draw.Point{X: 100, Y: 100},
    draw.Point{X: 200, Y: 200},
    draw.WithColor(draw.ColorBlue),
    draw.WithLineWidth(3.0),
)

// 创建填充矩形
filledRect := draw.NewRectangle(
    draw.Point{X: 300, Y: 150},
    draw.Point{X: 400, Y: 250},
    draw.WithColor(draw.ColorRed),
    draw.WithFill(true),
)

// 创建圆形
circle := draw.NewCircle(
    draw.Point{X: 500, Y: 300},
    50,
    draw.WithColor(draw.ColorGreen),
    draw.WithLineWidth(2.5),
)

// 创建图像处理器并添加图形
processor := draw.NewImageProcessor(
    "input.jpg",  // 也支持在线URL
    draw.WithOutputDir("result"),
    draw.WithTimeBasedName(),
    draw.WithShape(rect),
    draw.WithShape(filledRect),
    draw.WithShape(circle),
)

// 处理图像
if err := processor.Process(); err != nil {
    log.Fatalf("处理图像失败: %v", err)
}
```

### 使用工厂模式创建图形

```go
// 通过类型名称创建图形
circle, ok := draw.NewShape("circle",
    draw.WithPoints([]draw.Point{{X: 400, Y: 400}}),
    draw.WithRadius(60),
    draw.WithColor(draw.ColorRed),
    draw.WithFill(true),
)

if !ok {
    log.Fatal("创建圆形失败")
}
```

### 设置输出格式和质量

```go
processor := draw.NewImageProcessor(
    "input.jpg",
    draw.WithOutputDir("jpeg_results"),
    draw.WithOutputFormat(draw.FormatJPEG),
    draw.WithJpegQuality(75),
    draw.WithTimeBasedName(),
)
```

### 使用预处理和后处理功能

```go
processor := draw.NewImageProcessor(
    "input.jpg",
    draw.WithOutputDir("processed_results"),
    // 预处理函数 - 在绘制形状前应用灰度效果
    draw.WithPreProcess(func(dc *gg.Context, width, height float64) error {
        // 图像处理逻辑
        return nil
    }),
    // 后处理函数 - 添加文本水印
    draw.WithPostProcess(func(dc *gg.Context, width, height float64) error {
        dc.SetRGB(1, 1, 1) // 白色
        dc.DrawStringAnchored("© 水印", width-10, height-10, 1, 1)
        return nil
    }),
)
```

## 高级功能

### 自定义图形

可以通过实现Shape接口创建自定义图形：

```go
type MyShape struct {
    base.BaseShape
    // 自定义字段
}

func (s *MyShape) Draw(dc *gg.Context, width, height float64) error {
    // 实现绘制逻辑
    return nil
}

// 创建工厂
type MyShapeFactory struct{}

func (f MyShapeFactory) Create(options ...base.Option) base.Shape {
    shape := &MyShape{
        BaseShape: base.BaseShape{
            ShapeType: "my_shape",
            Color:     base.ColorBlack,
            LineWidth: 1.0,
        },
    }
    
    // 应用选项
    for _, option := range options {
        option(shape)
    }
    
    return shape
}

// 注册自定义图形
draw.RegisterShape("my_shape", MyShapeFactory{})
```

### 创建图像处理器

```go
processor := draw.NewImageProcessor(imagePath, options...)
```

### 图像处理器选项

```go
// 设置输出文件名
draw.WithOutputName("result.png")

// 使用时间戳作为文件名
draw.WithTimeBasedName()

// 设置输出目录
draw.WithOutputDir("outputs")

// 添加一条线
draw.WithLine(line)

// 添加多条线
draw.WithLines([]draw.Line{line1, line2})
```

### 处理图像

```go
err := processor.Process()
```

## 预设颜色

该库提供了以下预设颜色：

- `ColorWhite`: 白色 [1, 1, 1]
- `ColorBlack`: 黑色 [0, 0, 0]
- `ColorRed`: 红色 [1, 0, 0]
- `ColorBlue`: 蓝色 [0, 0, 1]
- `ColorGreen`: 绿色 [0, 1, 0]
- `ColorYellow`: 黄色 [1, 1, 0]

## 协议

MIT License 