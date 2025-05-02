package processor

import (
	"time"

	"github.com/fogleman/gg"
	"github.com/jiushengTech/common/utils/draw/shape/base"
)

// ProcessFunc 定义图像处理函数类型
type ProcessFunc func(dc *gg.Context, width, height float64) error

// ImageProcessor 图像处理器
type ImageProcessor struct {
	Path           string        // 图像文件路径
	Shapes         []base.Shape  // 图形集合
	Output         string        // 输出文件名
	OutputDir      string        // 输出目录
	Format         OutputFormat  // 输出格式
	JpegQuality    int           // JPEG质量 (1-100)
	RequestTimeout time.Duration // HTTP请求超时
	PreProcess     ProcessFunc   // 预处理函数
	PostProcess    ProcessFunc   // 后处理函数
}
