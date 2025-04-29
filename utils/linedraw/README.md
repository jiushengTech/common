# LineDrawer 线条绘制工具库

这是一个可用于在图像上绘制线条和数值的Go语言工具库。它支持在图像上绘制水平线、垂直线，以及它们之间的数值标注。

## 功能特点

- 在图像上绘制水平线和垂直线
- 在线条之间显示数值
- 自定义线条颜色和粗细
- 支持自动按坐标排序
- 基于选项模式的API设计
- 支持使用时间戳生成文件名

## 安装

```bash
go get github.com/yourusername/linedraw
```

## 快速开始

以下是一个简单的示例，展示如何使用该库在图像上绘制垂直线和水平线：

```go
package main

import (
    "fmt"
    "path/filepath"

    "github.com/yourusername/linedraw"
)

func main() {
    // 图片路径
    imagePath := filepath.Join("resource", "source.jpg")

    // 创建垂直线点和值
    xpoints := []linedraw.Point{
        {X: 279, Y: 0},
        {X: 380, Y: 0},
        {X: 494, Y: 0},
    }
    xvalues := []float64{0.33, 0.45}

    // 使用选项模式创建垂直线
    vline := linedraw.NewVerticalLine(
        xpoints, 
        xvalues, 
        linedraw.WithColor(linedraw.ColorWhite),
        linedraw.WithLineWidth(2.5),
    )

    // 创建图像处理器
    processor := linedraw.NewImageProcessor(
        imagePath,
        linedraw.WithTimeBasedName(),
        linedraw.WithLine(vline),
    )

    // 处理图像
    if err := processor.Process(); err != nil {
        fmt.Printf("错误: %v\n", err)
    } else {
        fmt.Printf("成功生成图片: %s/%s\n", processor.OutputDir, processor.Output)
    }
}
```

## API 文档

### 创建线条

```go
// 创建垂直线
line := linedraw.NewVerticalLine(points, values, options...)

// 创建水平线
line := linedraw.NewHorizontalLine(points, values, options...)
```

### 线条选项

```go
// 设置线条颜色
linedraw.WithColor(linedraw.ColorRed)

// 设置线条宽度
linedraw.WithLineWidth(2.5)
```

### 创建图像处理器

```go
processor := linedraw.NewImageProcessor(imagePath, options...)
```

### 图像处理器选项

```go
// 设置输出文件名
linedraw.WithOutputName("result.png")

// 使用时间戳作为文件名
linedraw.WithTimeBasedName()

// 设置输出目录
linedraw.WithOutputDir("outputs")

// 添加一条线
linedraw.WithLine(line)

// 添加多条线
linedraw.WithLines([]linedraw.Line{line1, line2})
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