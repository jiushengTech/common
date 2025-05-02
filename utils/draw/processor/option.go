package processor

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/jiushengTech/common/utils/draw/shape/base"
)

// Option 是图像设置的函数选项接口
type Option func(*ImageProcessor)

// WithOutputName 设置输出文件名
func WithOutputName(name string) Option {
	return func(p *ImageProcessor) {
		if name == "" {
			p.Output = GetTimeBasedFileName(p.Format)
		} else {
			// 确保扩展名与格式一致
			ext := string(p.Format)
			if !strings.HasSuffix(name, "."+ext) && ext != "" {
				baseName := strings.TrimSuffix(name, filepath.Ext(name))
				name = baseName + "." + ext
			}
			p.Output = name
		}
	}
}

// WithTimeBasedName 设置输出文件名为当前时间
func WithTimeBasedName() Option {
	return func(p *ImageProcessor) {
		p.Output = GetTimeBasedFileName(p.Format)
	}
}

// WithOutputDir 设置输出目录
func WithOutputDir(dir string) Option {
	return func(p *ImageProcessor) {
		if dir == "" {
			dir = "result"
		}
		p.OutputDir = dir
	}
}

// WithShape 添加一个图形
func WithShape(s base.Shape) Option {
	return func(p *ImageProcessor) {
		p.Shapes = append(p.Shapes, s)
	}
}

// WithShapes 添加多个图形
func WithShapes(shapes []base.Shape) Option {
	return func(p *ImageProcessor) {
		p.Shapes = append(p.Shapes, shapes...)
	}
}

// WithOutputFormat 设置输出格式
func WithOutputFormat(format OutputFormat) Option {
	return func(p *ImageProcessor) {
		p.Format = format

		// 自动更新文件扩展名
		ext := string(format)
		if !strings.HasSuffix(p.Output, "."+ext) && ext != "" {
			baseName := strings.TrimSuffix(p.Output, filepath.Ext(p.Output))
			p.Output = baseName + "." + ext
		}
	}
}

// WithJpegQuality 设置JPEG质量
func WithJpegQuality(quality int) Option {
	return func(p *ImageProcessor) {
		if quality < 1 {
			quality = 1
		} else if quality > 100 {
			quality = 100
		}
		p.JpegQuality = quality
	}
}

// WithRequestTimeout 设置HTTP请求超时时间
func WithRequestTimeout(timeout time.Duration) Option {
	return func(p *ImageProcessor) {
		if timeout > 0 {
			p.RequestTimeout = timeout
		}
	}
}

// WithPreProcess 设置预处理函数
func WithPreProcess(fn ProcessFunc) Option {
	return func(p *ImageProcessor) {
		p.PreProcess = fn
	}
}

// WithPostProcess 设置后处理函数
func WithPostProcess(fn ProcessFunc) Option {
	return func(p *ImageProcessor) {
		p.PostProcess = fn
	}
}

// GetTimeBasedFileName 返回以当前时间格式化的文件名
func GetTimeBasedFileName(format OutputFormat) string {
	now := time.Now()
	ext := "png"
	if format != "" {
		ext = string(format)
	}

	return fmt.Sprintf("%d%02d%02d_%02d%02d%02d.%s",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		ext)
}
