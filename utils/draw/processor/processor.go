package processor

import (
	"fmt"
	"image/jpeg"
	"io"
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

var (
	// 用于管理临时文件的清理
	tempFilesMutex sync.Mutex
	tempFiles      = make(map[string]bool)
	tempDir        = "temp_images" // 固定临时目录名
)

// 默认输出文件名
const DefaultOutputName = "output.png"

// 下载图片函数 - 内部使用
func downloadImage(url string, timeout time.Duration) (string, error) {
	// 设置默认超时
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	// 创建临时目录
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %v", err)
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

	// 下载图片 - 设置超时
	client := &http.Client{
		Timeout: timeout,
	}

	// 创建请求以便添加头部
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 添加用户代理
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("下载图片失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载图片失败，HTTP状态码: %d", resp.StatusCode)
	}

	// 创建本地文件
	out, err := os.Create(localPath)
	if err != nil {
		return "", fmt.Errorf("创建本地文件失败: %v", err)
	}
	defer out.Close()

	// 写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		// 删除失败的临时文件
		os.Remove(localPath)
		return "", fmt.Errorf("保存图片失败: %v", err)
	}

	return localPath, nil
}

// 清理临时文件
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

// 清理临时目录
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

// NewImageProcessor 创建一个新的图像处理器
func NewImageProcessor(imagePath string, options ...Option) *ImageProcessor {
	processor := &ImageProcessor{
		Path:           imagePath,
		Output:         DefaultOutputName,
		OutputDir:      "result",
		Format:         FormatPNG,
		JpegQuality:    90,
		RequestTimeout: 30 * time.Second,
		PreProcess:     nil,
		PostProcess:    nil,
	}

	// 应用所有选项
	for _, option := range options {
		option(processor)
	}

	return processor
}

// AddShape 添加一个图形
func (p *ImageProcessor) AddShape(s base.Shape) *ImageProcessor {
	p.Shapes = append(p.Shapes, s)
	return p
}

// AddShapes 添加多个图形
func (p *ImageProcessor) AddShapes(shapes []base.Shape) *ImageProcessor {
	p.Shapes = append(p.Shapes, shapes...)
	return p
}

// SetOutputFormat 设置输出格式
func (p *ImageProcessor) SetOutputFormat(format OutputFormat) *ImageProcessor {
	p.Format = format

	// 自动更新文件扩展名
	ext := string(format)
	if !strings.HasSuffix(p.Output, "."+ext) {
		base := strings.TrimSuffix(p.Output, filepath.Ext(p.Output))
		p.Output = base + "." + ext
	}

	return p
}

// SetJpegQuality 设置JPEG质量
func (p *ImageProcessor) SetJpegQuality(quality int) *ImageProcessor {
	if quality < 1 {
		quality = 1
	} else if quality > 100 {
		quality = 100
	}
	p.JpegQuality = quality
	return p
}

// Validate 验证图像和图形数据是否有效
func (p *ImageProcessor) Validate() error {
	// 检查图像路径
	if p.Path == "" {
		return fmt.Errorf("图像路径为空")
	}

	// 验证输出格式
	if p.Format != FormatPNG && p.Format != FormatJPEG {
		return fmt.Errorf("不支持的输出格式: %s", p.Format)
	}

	return nil
}

// Process 处理图像并绘制图形，返回图片的绝对路径
func (p *ImageProcessor) Process() (string, error) {
	// 验证数据
	if err := p.Validate(); err != nil {
		return "", err
	}

	// 处理图片路径，如果是URL则下载
	imagePath := p.Path
	isTemporaryFile := false
	if strings.HasPrefix(p.Path, "http://") || strings.HasPrefix(p.Path, "https://") {
		localPath, err := downloadImage(p.Path, p.RequestTimeout)
		if err != nil {
			return "", fmt.Errorf("下载图片失败: %v", err)
		}
		imagePath = localPath
		isTemporaryFile = true

		// 确保函数结束时删除临时文件
		defer func() {
			if isTemporaryFile {
				cleanupTempFile(imagePath)
			}
		}()
	}

	// 加载原图
	sourceImage, err := gg.LoadImage(imagePath)
	if err != nil {
		return "", fmt.Errorf("加载图像 %s 失败: %v", imagePath, err)
	}

	// 获取原图尺寸
	bounds := sourceImage.Bounds()
	width := float64(bounds.Max.X)
	height := float64(bounds.Max.Y)

	// 创建新的画布
	dc := gg.NewContext(bounds.Max.X, bounds.Max.Y)

	// 绘制原图
	dc.DrawImage(sourceImage, 0, 0)

	// 执行预处理函数（如果有）
	if p.PreProcess != nil {
		if err := p.PreProcess(dc, width, height); err != nil {
			return "", fmt.Errorf("图像预处理失败: %v", err)
		}
	}

	// 绘制每个图形
	if err := p.drawShapes(dc, width, height); err != nil {
		return "", err
	}

	// 执行后处理函数（如果有）
	if p.PostProcess != nil {
		if err := p.PostProcess(dc, width, height); err != nil {
			return "", fmt.Errorf("图像后处理失败: %v", err)
		}
	}

	// 创建输出目录（如果不存在）
	if err := os.MkdirAll(p.OutputDir, 0755); err != nil {
		return "", fmt.Errorf("创建目录 %s 失败: %v", p.OutputDir, err)
	}

	// 保存结果
	outputPath := filepath.Join(p.OutputDir, p.Output)

	// 根据格式保存文件
	var saveErr error
	switch p.Format {
	case FormatPNG:
		saveErr = dc.SavePNG(outputPath)
	case FormatJPEG:
		// 创建文件
		file, err := os.Create(outputPath)
		if err != nil {
			return "", fmt.Errorf("创建JPEG文件失败: %v", err)
		}
		defer file.Close()

		// 使用jpeg库保存图像
		img := dc.Image()
		opt := jpeg.Options{Quality: p.JpegQuality}
		saveErr = jpeg.Encode(file, img, &opt)
		if saveErr != nil {
			saveErr = fmt.Errorf("保存JPEG文件失败: %v", saveErr)
		}
	default:
		return "", fmt.Errorf("不支持的输出格式: %s", p.Format)
	}

	if saveErr != nil {
		return "", saveErr
	}

	// 获取绝对路径并返回
	absPath, err := p.GetAbsoluteOutputPath()
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// GetOutputPath 获取完整的输出路径
func (p *ImageProcessor) GetOutputPath() string {
	return filepath.Join(p.OutputDir, p.Output)
}

// GetAbsoluteOutputPath 获取图片输出的绝对路径
func (p *ImageProcessor) GetAbsoluteOutputPath() (string, error) {
	relPath := p.GetOutputPath()
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return "", fmt.Errorf("获取绝对路径失败: %v", err)
	}
	return absPath, nil
}

// drawShapes 绘制所有图形
func (p *ImageProcessor) drawShapes(dc *gg.Context, width, height float64) error {
	for i, s := range p.Shapes {
		if err := s.Draw(dc, width, height); err != nil {
			return fmt.Errorf("绘制图形 %d 失败: %v", i, err)
		}
	}
	return nil
}
