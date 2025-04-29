package linedraw

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"golang.org/x/image/font/basicfont"
)

// ImageProcessor 图像处理器
type ImageProcessor struct {
	Path      string // 图像文件路径
	Lines     []Line // 线条集合
	Output    string // 输出文件名
	OutputDir string // 输出目录
}

// NewImageProcessor 创建一个新的图像处理器
func NewImageProcessor(imagePath string, options ...ImageOption) *ImageProcessor {
	processor := &ImageProcessor{
		Path:      imagePath,
		Output:    DefaultOutputName,
		OutputDir: "result",
	}

	// 应用所有选项
	for _, option := range options {
		option(processor)
	}

	return processor
}

// AddLine 添加一条线
func (p *ImageProcessor) AddLine(line Line) *ImageProcessor {
	p.Lines = append(p.Lines, line)
	return p
}

// Validate 验证图像和线条数据是否有效
func (p *ImageProcessor) Validate() error {
	// 检查图像路径
	if p.Path == "" {
		return fmt.Errorf("image path is empty")
	}

	// 检查每条线
	for i, line := range p.Lines {
		if len(line.Points) < 2 {
			return fmt.Errorf("line %d must have at least 2 points", i)
		}
		if len(line.Values) != len(line.Points)-1 {
			return fmt.Errorf("line %d values count must be points count - 1", i)
		}
	}
	return nil
}

// Process 处理图像并绘制线条
func (p *ImageProcessor) Process() error {
	// 验证数据
	if err := p.Validate(); err != nil {
		return err
	}

	// 加载原图
	sourceImage, err := gg.LoadImage(p.Path)
	if err != nil {
		return fmt.Errorf("failed to load image %s: %v", p.Path, err)
	}

	// 获取原图尺寸
	bounds := sourceImage.Bounds()
	width := float64(bounds.Max.X)
	height := float64(bounds.Max.Y)

	// 创建新的画布
	dc := gg.NewContext(bounds.Max.X, bounds.Max.Y)

	// 绘制原图
	dc.DrawImage(sourceImage, 0, 0)

	// 绘制每条线及其值
	for _, line := range p.Lines {
		if line.Type == VerticalLine {
			// 绘制垂直线
			for i, point := range line.Points {
				// 验证坐标
				if point.X < 0 || point.X > width {
					return fmt.Errorf("x coordinate %.2f is out of range [0, %.2f]", point.X, width)
				}

				// 画线
				dc.SetRGB(line.Color[0], line.Color[1], line.Color[2])
				dc.SetLineWidth(line.LineWidth)
				dc.DrawLine(point.X, 0, point.X, height)
				dc.Stroke()

				// 如果不是最后一个点，绘制值
				if i < len(line.Points)-1 {
					// 计算文字位置（两条线的中间位置）
					textX := (point.X + line.Points[i+1].X) / 2
					textY := height / 3 // 在1/3处显示文字

					// 绘制文字
					text := fmt.Sprintf("%.2f", line.Values[i])
					drawText(dc, text, textX, textY)
				}
			}
		} else if line.Type == HorizontalLine {
			// 绘制水平线
			for i, point := range line.Points {
				// 验证坐标
				if point.Y < 0 || point.Y > height {
					return fmt.Errorf("y coordinate %.2f is out of range [0, %.2f]", point.Y, height)
				}

				// 画线
				dc.SetRGB(line.Color[0], line.Color[1], line.Color[2])
				dc.SetLineWidth(line.LineWidth)
				dc.DrawLine(0, point.Y, width, point.Y)
				dc.Stroke()

				// 如果不是最后一个点，绘制值
				if i < len(line.Points)-1 {
					// 计算文字位置（两条线的中间位置）
					textY := (point.Y + line.Points[i+1].Y) / 2
					textX := width / 3 // 在1/3处显示文字

					// 绘制文字
					text := fmt.Sprintf("%.2f", line.Values[i])
					drawText(dc, text, textX, textY)
				}
			}
		}
	}

	// 创建输出目录（如果不存在）
	if err := os.MkdirAll(p.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", p.OutputDir, err)
	}

	// 保存结果
	outputPath := filepath.Join(p.OutputDir, p.Output)
	return dc.SavePNG(outputPath)
}

// drawText 绘制文本（带描边效果）
func drawText(dc *gg.Context, text string, x, y float64) {
	// 设置文字
	face := basicfont.Face7x13
	dc.SetFontFace(face)

	// 绘制黑色描边
	dc.SetLineWidth(3)
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(text, x, y, 0.5, 0.5)
	dc.Stroke()

	// 绘制白色文字
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(text, x, y, 0.5, 0.5)
	dc.Fill()
}
