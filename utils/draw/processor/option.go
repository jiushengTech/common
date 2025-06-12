package processor

import (
	"github.com/jiushengTech/common/utils/draw/shape/base"
	"path/filepath"
	"strings"
)

// Option 是图像设置的函数选项接口
type Option func(*ImageProcessor)

// WithOutputName 设置输出文件名
func WithOutputName(name string) Option {
	return func(p *ImageProcessor) {
		if name == "" {
			p.OutputName = GetDefaultOutputName(p.Format)
		} else {
			// 确保扩展名与格式一致
			ext := string(p.Format)
			if !strings.HasSuffix(name, "."+ext) && ext != "" {
				baseName := strings.TrimSuffix(name, filepath.Ext(name))
				name = baseName + "." + ext
			}
			p.OutputName = name
		}
	}
}

// WithTimeBasedName 设置输出文件名为当前时间
func WithTimeBasedName() Option {
	return func(p *ImageProcessor) {
		p.OutputName = GetDefaultOutputName(p.Format)
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
		if !strings.HasSuffix(p.OutputName, "."+ext) && ext != "" {
			baseName := strings.TrimSuffix(p.OutputName, filepath.Ext(p.OutputName))
			p.OutputName = baseName + "." + ext
		}
	}
}

// WithJpegQuality 设置JPEG质量
func WithJpegQuality(quality int) Option {
	return func(p *ImageProcessor) {
		// 限制值在0-100范围内
		if quality < 0 {
			quality = 0
		} else if quality > 100 {
			quality = 100
		}
		p.JpegQuality = quality
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
