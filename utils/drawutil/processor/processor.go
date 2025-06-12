package processor

import (
	"fmt"
	"image"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jiushengTech/common/utils/fileutil"

	"github.com/jiushengTech/common/utils/drawutil/shape/base"

	"github.com/fogleman/gg"
)

// OutputFormat 表示输出图像格式
type OutputFormat string

// 支持的输出格式
const (
	FormatPNG  OutputFormat = "png"  // PNG格式
	FormatJPEG OutputFormat = "jpeg" // JPEG格式
)

// ProcessFunc 定义图像处理函数类型
type ProcessFunc func(dc *gg.Context, width, height float64) error

// ImageProcessor 图像处理器，管理图像处理流程
type ImageProcessor struct {
	// 输入配置
	Path string // 图像路径或URL
	// 输出配置
	OutputName  string       // 输出文件名，不包含路径
	OutputDir   string       // 输出目录
	Format      OutputFormat // 输出图片格式
	JpegQuality int          // JPEG质量 (0-100)

	// 处理配置
	Shapes      []base.Shape // 要绘制的图形集合
	PreProcess  ProcessFunc  // 预处理函数
	PostProcess ProcessFunc  // 后处理函数
}

// NewImageProcessor 创建一个新的图像处理器
func NewImageProcessor(imagePath string, options ...Option) *ImageProcessor {
	processor := &ImageProcessor{
		Path:        imagePath,
		OutputName:  GetDefaultOutputName(FormatPNG), // 使用基于时间的默认文件名
		Format:      FormatPNG,
		JpegQuality: 100,
		PreProcess:  nil,
		PostProcess: nil,
		Shapes:      []base.Shape{},
	}

	// 应用所有选项
	for _, option := range options {
		option(processor)
	}

	return processor
}

// Process 处理图像并保存到输出路径（优化版）
func (p *ImageProcessor) Process() (string, error) {
	// 验证输入路径
	if p.Path == "" {
		return "", fmt.Errorf("输入路径不能为空")
	}

	// 判断是否为URL
	isURL := strings.HasPrefix(strings.ToLower(p.Path), "http://") ||
		strings.HasPrefix(strings.ToLower(p.Path), "https://")

	// 验证本地文件是否存在
	if !isURL {
		if _, err := os.Stat(p.Path); os.IsNotExist(err) {
			return "", fmt.Errorf("输入文件不存在: %s", p.Path)
		}
	}

	// 验证输出格式
	if p.Format != FormatPNG && p.Format != FormatJPEG {
		return "", fmt.Errorf("不支持的输出格式: %s", p.Format)
	}

	// 处理图像
	img, err := p.loadImage()
	if err != nil {
		return "", err
	}

	// 创建图形上下文
	width := float64(img.Bounds().Dx())
	height := float64(img.Bounds().Dy())
	dc := gg.NewContextForImage(img)

	// 应用预处理函数
	if p.PreProcess != nil {
		if err := p.PreProcess(dc, width, height); err != nil {
			return "", fmt.Errorf("预处理失败: %w", err)
		}
	}

	// 绘制所有图形（可以考虑并行处理）
	if err := p.drawShapes(dc, width, height); err != nil {
		return "", err
	}

	// 应用后处理函数
	if p.PostProcess != nil {
		if err := p.PostProcess(dc, width, height); err != nil {
			return "", fmt.Errorf("后处理失败: %w", err)
		}
	}

	// 保存结果
	return p.saveImage(dc)
}

// loadImage 加载图像（从URL或本地文件）
func (p *ImageProcessor) loadImage() (image.Image, error) {
	// 判断是否为URL
	isURL := strings.HasPrefix(strings.ToLower(p.Path), "http://") ||
		strings.HasPrefix(strings.ToLower(p.Path), "https://")

	// 如果是URL，下载图片
	imagePath := p.Path
	if isURL {
		var err error
		imagePath, err = fileutil.DownloadFile(p.Path, ".")
		if err != nil {
			return nil, fmt.Errorf("下载图片失败: %w", err)
		}
		// 在函数结束时清理临时文件
		defer fileutil.CleanUp(imagePath)
	}

	// 加载图片
	img, err := gg.LoadImage(imagePath)
	if err != nil {
		return nil, fmt.Errorf("加载图片失败 [%s]: %w", imagePath, err)
	}

	return img, nil
}

// saveImage 保存处理后的图像
func (p *ImageProcessor) saveImage(dc *gg.Context) (string, error) {
	// 确保输出目录存在
	if p.OutputDir == "" {
		p.OutputDir = "."
	}
	if err := os.MkdirAll(p.OutputDir, 0755); err != nil {
		return "", fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 获取输出路径
	outputPath := filepath.Join(p.OutputDir, p.OutputName)

	// 保存结果图像
	switch p.Format {
	case FormatPNG:
		if err := dc.SavePNG(outputPath); err != nil {
			return "", fmt.Errorf("保存PNG图像失败: %w", err)
		}
	case FormatJPEG:
		file, err := os.Create(outputPath)
		if err != nil {
			return "", fmt.Errorf("创建JPEG输出文件失败: %w", err)
		}
		defer file.Close()

		if err := jpeg.Encode(file, dc.Image(), &jpeg.Options{
			Quality: p.JpegQuality,
		}); err != nil {
			return "", fmt.Errorf("保存JPEG图像失败: %w", err)
		}
	default:
		return "", fmt.Errorf("不支持的输出格式: %s", p.Format)
	}

	// 返回绝对路径
	return p.GetAbsoluteOutputPath()
}

// 优化的drawShapes方法，对于大量图形时可以考虑并行处理
func (p *ImageProcessor) drawShapes(dc *gg.Context, width, height float64) error {
	// 如果图形数量较多，考虑并行处理
	if len(p.Shapes) > 10 {
		var wg sync.WaitGroup
		errors := make(chan error, len(p.Shapes))

		// 创建一个临时上下文的副本，用于后续合并
		tempContexts := make([]*gg.Context, len(p.Shapes))

		for i, shape := range p.Shapes {
			wg.Add(1)
			go func(idx int, s base.Shape) {
				defer wg.Done()
				// 为每个图形创建新的上下文
				tempDC := gg.NewContext(int(width), int(height))
				if err := s.Draw(tempDC, width, height); err != nil {
					errors <- fmt.Errorf("绘制图形 %s 失败: %w", s.GetType(), err)
					return
				}
				tempContexts[idx] = tempDC
			}(i, shape)
		}

		// 等待所有绘制任务完成
		wg.Wait()
		close(errors)

		// 检查是否有错误
		for err := range errors {
			if err != nil {
				return err
			}
		}

		// 合并所有上下文
		for _, tempDC := range tempContexts {
			if tempDC != nil {
				dc.DrawImage(tempDC.Image(), 0, 0)
			}
		}

		return nil
	}

	// 如果图形数量较少，直接串行处理
	for _, shape := range p.Shapes {
		if err := shape.Draw(dc, width, height); err != nil {
			return fmt.Errorf("绘制图形 %s 失败: %w", shape.GetType(), err)
		}
	}
	return nil
}

// 添加一个调整大小的函数
func (p *ImageProcessor) ResizeImage(maxWidth, maxHeight int) Option {
	return func(p *ImageProcessor) {
		oldPreProcess := p.PreProcess

		p.PreProcess = func(dc *gg.Context, width, height float64) error {
			// 先执行原来的预处理
			if oldPreProcess != nil {
				if err := oldPreProcess(dc, width, height); err != nil {
					return err
				}
			}

			// 计算新尺寸，保持宽高比
			newWidth, newHeight := calculateNewSize(int(width), int(height), maxWidth, maxHeight)

			// 创建新的图形上下文并调整大小
			newDC := gg.NewContext(newWidth, newHeight)
			newDC.DrawImage(dc.Image(), 0, 0)

			// 替换原始上下文
			*dc = *newDC

			return nil
		}
	}
}

// 计算调整后的尺寸，保持宽高比
func calculateNewSize(width, height, maxWidth, maxHeight int) (int, int) {
	if width <= maxWidth && height <= maxHeight {
		return width, height
	}

	widthRatio := float64(maxWidth) / float64(width)
	heightRatio := float64(maxHeight) / float64(height)

	ratio := math.Min(widthRatio, heightRatio)

	return int(float64(width) * ratio), int(float64(height) * ratio)
}

// GetOutputPath 获取输出路径
func (p *ImageProcessor) GetOutputPath() string {
	return filepath.Join(p.OutputDir, p.OutputName)
}

// GetAbsoluteOutputPath 获取绝对输出路径
func (p *ImageProcessor) GetAbsoluteOutputPath() (string, error) {
	outputPath := p.GetOutputPath()
	absolutePath, err := filepath.Abs(outputPath)
	if err != nil {
		return "", fmt.Errorf("获取绝对路径失败: %w", err)
	}
	return absolutePath, nil
}

// GetDefaultOutputName 获取默认输出文件名，使用当前时间格式（毫秒级）
func GetDefaultOutputName(format OutputFormat) string {
	// 获取当前时间
	now := time.Now()
	// 默认扩展名为 png
	ext := "png" // 默认格式
	if format != "" {
		ext = string(format) // 如果传入了格式，使用传入的格式
	}

	// 获取毫秒级时间戳
	milliTimestamp := now.UnixMilli() // 毫秒级时间戳

	// 格式化文件名：使用年月日时分秒和毫秒来确保唯一性
	return fmt.Sprintf("%d%02d%02d_%02d%02d%02d_%03d.%s",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		milliTimestamp%1000, ext) // 使用毫秒的后3位
}
