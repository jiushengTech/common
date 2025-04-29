package linedraw

import (
	"fmt"
	"time"
)

// ImageOption 是图像设置的函数选项接口
type ImageOption func(*ImageProcessor)

// WithOutputName 设置输出文件名
func WithOutputName(name string) ImageOption {
	return func(p *ImageProcessor) {
		if name == "" {
			p.Output = GetTimeBasedFileName()
		} else {
			p.Output = name
		}
	}
}

// WithTimeBasedName 设置输出文件名为当前时间
func WithTimeBasedName() ImageOption {
	return func(p *ImageProcessor) {
		p.Output = GetTimeBasedFileName()
	}
}

// WithOutputDir 设置输出目录
func WithOutputDir(dir string) ImageOption {
	return func(p *ImageProcessor) {
		p.OutputDir = dir
	}
}

// WithLine 添加一条线
func WithLine(line Line) ImageOption {
	return func(p *ImageProcessor) {
		p.Lines = append(p.Lines, line)
	}
}

// WithLines 添加多条线
func WithLines(lines []Line) ImageOption {
	return func(p *ImageProcessor) {
		p.Lines = append(p.Lines, lines...)
	}
}

// GetTimeBasedFileName 返回以当前时间格式化的文件名
func GetTimeBasedFileName() string {
	now := time.Now()
	return fmt.Sprintf("%d%02d%02d_%02d%02d%02d.png",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
}

// LineOption 是线条设置的函数选项接口
type LineOption func(*Line)

// WithColor 设置线条颜色
func WithColor(color [3]float64) LineOption {
	return func(l *Line) {
		l.Color = color
	}
}

// WithLineWidth 设置线条宽度
func WithLineWidth(width float64) LineOption {
	return func(l *Line) {
		l.LineWidth = width
	}
}
