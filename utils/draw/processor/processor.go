package processor

import (
	"context"
	"fmt"
	"image/jpeg"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jiushengTech/common/utils/draw/shape/base"

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

var (
	// 用于管理临时文件的清理
	tempFilesMutex sync.Mutex
	tempFiles      = make(map[string]bool)
	tempDir        = "temp_images" // 固定临时目录名
)

// GetDefaultOutputName 获取默认输出文件名，使用当前时间格式
func GetDefaultOutputName(format OutputFormat) string {
	now := time.Now()
	ext := "png"
	if format != "" {
		ext = string(format)
	}
	rand.Seed(now.UnixNano())      // 设置随机种子
	randomID := rand.Intn(1000000) // 生成 0~999999 的随机数
	return fmt.Sprintf("%d%02d%02d_%02d%02d%02d_%06d.%s",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		randomID, ext)

}

// 默认输出文件名
const DefaultOutputName = "output.png" // 兼容旧代码，实际将使用时间格式

// downloadImage 下载图片函数 - 内部使用
func downloadImage(url string, timeout time.Duration) (string, error) {
	// 设置默认超时
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	// 创建临时目录
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	// 生成本地文件名 - 使用时间戳和URL哈希避免文件名冲突
	timestamp := time.Now().UnixNano()
	urlBase := filepath.Base(url)
	fileName := fmt.Sprintf("%d_%s", timestamp, urlBase)

	// 处理URL没有有效文件名的情况
	if fileName == fmt.Sprintf("%d_", timestamp) || fileName == fmt.Sprintf("%d_.", timestamp) {
		fileName = fmt.Sprintf("%d_temp_image.jpg", timestamp)
	}
	localPath := filepath.Join(tempDir, fileName)

	// 添加到临时文件列表，以便后续清理
	tempFilesMutex.Lock()
	tempFiles[localPath] = true
	tempFilesMutex.Unlock()

	// 使用context设置超时
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 创建请求以便添加头部
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 添加用户代理
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	// 执行请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("下载图片失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载图片失败，HTTP状态码: %d", resp.StatusCode)
	}

	// 创建本地文件
	out, err := os.Create(localPath)
	if err != nil {
		return "", fmt.Errorf("创建本地文件失败: %w", err)
	}
	defer out.Close()

	// 写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		// 删除失败的临时文件
		os.Remove(localPath)
		return "", fmt.Errorf("保存图片失败: %w", err)
	}

	return localPath, nil
}

// cleanupTempFile 清理单个临时文件
func cleanupTempFile(path string) {
	tempFilesMutex.Lock()
	defer tempFilesMutex.Unlock()

	if _, exists := tempFiles[path]; exists {
		if err := os.Remove(path); err != nil {
			fmt.Printf("警告: 无法删除临时文件 %s: %v\n", path, err)
		} else {
			delete(tempFiles, path)

			// 如果没有更多临时文件，尝试删除目录
			if len(tempFiles) == 0 {
				cleanupTempDir()
			}
		}
	}
}

// cleanupTempDir 清理临时目录
func cleanupTempDir() {
	// 检查目录是否存在
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return // 目录不存在，无需删除
	}

	// 检查目录是否为空
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		fmt.Printf("警告: 无法读取临时目录 %s: %v\n", tempDir, err)
		return
	}

	// 如果目录不为空，不删除
	if len(entries) > 0 {
		return
	}

	// 删除空目录
	if err := os.Remove(tempDir); err != nil {
		fmt.Printf("警告: 无法删除临时目录 %s: %v\n", tempDir, err)
	}
}

// CleanupAllTempFiles 清理所有临时文件和目录
// 可以在程序结束时调用
func CleanupAllTempFiles() {
	tempFilesMutex.Lock()
	defer tempFilesMutex.Unlock()

	for path := range tempFiles {
		if err := os.Remove(path); err != nil {
			fmt.Printf("警告: 无法删除临时文件 %s: %v\n", path, err)
		} else {
			delete(tempFiles, path)
		}
	}

	cleanupTempDir()
}

// ImageProcessor 图像处理器，管理图像处理流程
type ImageProcessor struct {
	// 输入配置
	Path           string        // 图像路径或URL
	RequestTimeout time.Duration // HTTP请求超时时间

	// 输出配置
	Output      string       // 输出文件名，不包含路径
	OutputDir   string       // 输出目录
	Format      OutputFormat // 输出格式
	JpegQuality int          // JPEG质量 (0-100)

	// 处理配置
	Shapes      []base.Shape // 要绘制的图形集合
	PreProcess  ProcessFunc  // 预处理函数
	PostProcess ProcessFunc  // 后处理函数
}

// NewImageProcessor 创建一个新的图像处理器
func NewImageProcessor(imagePath string, options ...Option) *ImageProcessor {
	processor := &ImageProcessor{
		Path:           imagePath,
		Output:         GetDefaultOutputName(FormatPNG), // 使用基于时间的默认文件名
		OutputDir:      "result",
		Format:         FormatPNG,
		JpegQuality:    90,
		RequestTimeout: 30 * time.Second,
		PreProcess:     nil,
		PostProcess:    nil,
		Shapes:         []base.Shape{},
	}

	// 应用所有选项
	for _, option := range options {
		option(processor)
	}

	return processor
}

// AddShape 添加一个图形到处理器
func (p *ImageProcessor) AddShape(s base.Shape) *ImageProcessor {
	p.Shapes = append(p.Shapes, s)
	return p
}

// AddShapes 添加多个图形到处理器
func (p *ImageProcessor) AddShapes(shapes []base.Shape) *ImageProcessor {
	p.Shapes = append(p.Shapes, shapes...)
	return p
}

// SetOutputFormat 设置输出格式
func (p *ImageProcessor) SetOutputFormat(format OutputFormat) *ImageProcessor {
	p.Format = format

	// 如果文件名没有扩展名，或扩展名与格式不匹配，更新文件名扩展名
	if filepath.Ext(p.Output) == "" || !strings.HasSuffix(p.Output, string(format)) {
		baseName := strings.TrimSuffix(p.Output, filepath.Ext(p.Output))
		p.Output = baseName + "." + string(format)
	}

	return p
}

// SetJpegQuality 设置JPEG质量
func (p *ImageProcessor) SetJpegQuality(quality int) *ImageProcessor {
	// 限制值在0-100范围内
	if quality < 0 {
		quality = 0
	} else if quality > 100 {
		quality = 100
	}
	p.JpegQuality = quality
	return p
}

// Validate 验证处理器配置
func (p *ImageProcessor) Validate() error {
	// 验证输入路径
	if p.Path == "" {
		return fmt.Errorf("输入路径不能为空")
	}

	// 判断是否为URL
	isURL := strings.HasPrefix(strings.ToLower(p.Path), "http://") ||
		strings.HasPrefix(strings.ToLower(p.Path), "https://")

	// 验证本地文件是否存在
	if !isURL {
		if _, err := os.Stat(p.Path); os.IsNotExist(err) {
			return fmt.Errorf("输入文件不存在: %s", p.Path)
		}
	}

	// 验证输出格式
	if p.Format != FormatPNG && p.Format != FormatJPEG {
		return fmt.Errorf("不支持的输出格式: %s", p.Format)
	}

	return nil
}

// Process 处理图像并保存到输出路径
func (p *ImageProcessor) Process() (string, error) {
	// 验证配置
	if err := p.Validate(); err != nil {
		return "", err
	}

	// 判断是否为URL
	isURL := strings.HasPrefix(strings.ToLower(p.Path), "http://") ||
		strings.HasPrefix(strings.ToLower(p.Path), "https://")

	// 如果是URL，下载图片
	imagePath := p.Path
	if isURL {
		var err error
		imagePath, err = downloadImage(p.Path, p.RequestTimeout)
		if err != nil {
			return "", fmt.Errorf("下载图片失败: %w", err)
		}
		// 在函数结束时清理临时文件
		defer cleanupTempFile(imagePath)
	}

	// 加载图片
	img, err := gg.LoadImage(imagePath)
	if err != nil {
		return "", fmt.Errorf("加载图片失败: %w", err)
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

	// 绘制所有图形
	if err := p.drawShapes(dc, width, height); err != nil {
		return "", fmt.Errorf("绘制图形失败: %w", err)
	}

	// 应用后处理函数
	if p.PostProcess != nil {
		if err := p.PostProcess(dc, width, height); err != nil {
			return "", fmt.Errorf("后处理失败: %w", err)
		}
	}

	// 确保输出目录存在
	if err := os.MkdirAll(p.OutputDir, 0755); err != nil {
		return "", fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 获取输出路径
	outputPath := filepath.Join(p.OutputDir, p.Output)

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

// GetOutputPath 获取输出路径
func (p *ImageProcessor) GetOutputPath() string {
	return filepath.Join(p.OutputDir, p.Output)
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

// drawShapes 绘制所有图形
func (p *ImageProcessor) drawShapes(dc *gg.Context, width, height float64) error {
	for _, shape := range p.Shapes {
		if err := shape.Draw(dc, width, height); err != nil {
			return fmt.Errorf("绘制图形 %s 失败: %w", shape.GetType(), err)
		}
	}
	return nil
}
