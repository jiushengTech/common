package draw

import (
	"fmt"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

// Canvas 是对gg.Context的包装，提供更友好的绘图接口
type Canvas struct {
	dc     *gg.Context
	width  int
	height int
	font   *truetype.Font
}

// NewCanvas 创建一个新的画布
func NewCanvas(width, height int) *Canvas {
	// 创建gg.Context
	dc := gg.NewContext(width, height)

	// 加载默认字体
	f, _ := truetype.Parse(goregular.TTF)

	// 创建并返回画布
	return &Canvas{
		dc:     dc,
		width:  width,
		height: height,
		font:   f,
	}
}

// NewCanvasFromImage 从现有图像创建画布
func NewCanvasFromImage(img image.Image) *Canvas {
	// 创建gg.Context
	dc := gg.NewContextForImage(img)
	bounds := img.Bounds()

	// 加载默认字体
	f, _ := truetype.Parse(goregular.TTF)

	// 创建并返回画布
	return &Canvas{
		dc:     dc,
		width:  bounds.Dx(),
		height: bounds.Dy(),
		font:   f,
	}
}

// Width 返回画布宽度
func (c *Canvas) Width() int {
	return c.width
}

// Height 返回画布高度
func (c *Canvas) Height() int {
	return c.height
}

// SetColor 设置当前绘图颜色
func (c *Canvas) SetColor(color [3]float64) {
	c.dc.SetRGB(color[0], color[1], color[2])
}

// SetRGBA 设置带透明度的颜色
func (c *Canvas) SetRGBA(r, g, b, a float64) {
	c.dc.SetRGBA(r, g, b, a)
}

// SetLineWidth 设置线宽
func (c *Canvas) SetLineWidth(width float64) {
	c.dc.SetLineWidth(width)
}

// DrawLine 绘制线段
func (c *Canvas) DrawLine(x1, y1, x2, y2 float64) {
	c.dc.DrawLine(x1, y1, x2, y2)
	c.dc.Stroke()
}

// DrawRectangle 绘制矩形
func (c *Canvas) DrawRectangle(x, y, width, height float64, fill bool) {
	c.dc.DrawRectangle(x, y, width, height)
	if fill {
		c.dc.Fill()
	} else {
		c.dc.Stroke()
	}
}

// DrawCircle 绘制圆形
func (c *Canvas) DrawCircle(x, y, radius float64, fill bool) {
	c.dc.DrawCircle(x, y, radius)
	if fill {
		c.dc.Fill()
	} else {
		c.dc.Stroke()
	}
}

// DrawEllipse 绘制椭圆
func (c *Canvas) DrawEllipse(x, y, width, height float64, fill bool) {
	c.dc.DrawEllipse(x, y, width/2, height/2)
	if fill {
		c.dc.Fill()
	} else {
		c.dc.Stroke()
	}
}

// DrawText 绘制文本
func (c *Canvas) DrawText(text string, x, y float64) error {
	c.dc.DrawString(text, x, y)
	return nil
}

// DrawTextWithScale 绘制带缩放的文本
func (c *Canvas) DrawTextWithScale(text string, x, y, scaleX, scaleY float64) error {
	c.dc.Push()
	c.dc.Scale(scaleX, scaleY)
	c.dc.DrawString(text, x/scaleX, y/scaleY)
	c.dc.Pop()
	return nil
}

// SetFontFace 设置字体
func (c *Canvas) SetFontFace(size float64) {
	face := truetype.NewFace(c.font, &truetype.Options{
		Size: size,
	})
	c.dc.SetFontFace(face)
}

// LoadFont 加载自定义字体
func (c *Canvas) LoadFont(fontData []byte) error {
	font, err := truetype.Parse(fontData)
	if err != nil {
		return fmt.Errorf("加载字体失败: %w", err)
	}
	c.font = font
	return nil
}

// GetContext 获取底层的gg.Context实例
// 用于高级操作或直接访问底层功能
func (c *Canvas) GetContext() *gg.Context {
	return c.dc
}

// Image 返回生成的图像
func (c *Canvas) Image() image.Image {
	return c.dc.Image()
}

// Clear 清除画布
func (c *Canvas) Clear() {
	c.dc.Clear()
}

// ClearWithColor 使用指定颜色清除画布
func (c *Canvas) ClearWithColor(r, g, b float64) {
	c.dc.SetRGB(r, g, b)
	c.dc.Clear()
}

// SavePNG 保存为PNG图像
func (c *Canvas) SavePNG(path string) error {
	return c.dc.SavePNG(path)
}

// DrawImage 在指定位置绘制图像
func (c *Canvas) DrawImage(img image.Image, x, y float64) {
	c.dc.DrawImage(img, int(x), int(y))
}

// DrawImageScaled 绘制缩放后的图像
func (c *Canvas) DrawImageScaled(img image.Image, x, y, width, height float64) {
	// 缩放图像
	// 由于go标准库中没有内置的BiLinear，我们直接使用简单的缩放方法
	bounds := img.Bounds()
	scaledImg := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	// 使用简单的缩放算法
	scaleX := float64(bounds.Dx()) / width
	scaleY := float64(bounds.Dy()) / height

	for ny := 0; ny < int(height); ny++ {
		for nx := 0; nx < int(width); nx++ {
			// 计算原图坐标
			sx := int(float64(nx) * scaleX)
			sy := int(float64(ny) * scaleY)

			// 获取原图颜色并设置到新图
			c := img.At(bounds.Min.X+sx, bounds.Min.Y+sy)
			scaledImg.Set(nx, ny, c)
		}
	}

	// 绘制到画布
	c.dc.DrawImage(scaledImg, int(x), int(y))
}

// DrawImageWithOpacity 使用指定透明度绘制图像
func (c *Canvas) DrawImageWithOpacity(img image.Image, x, y float64, opacity float64) {
	c.dc.Push()
	c.dc.SetRGBA(1, 1, 1, opacity)
	c.dc.DrawImage(img, int(x), int(y))
	c.dc.Pop()
}

// Rotate 旋转画布（弧度制）
func (c *Canvas) Rotate(angle float64) {
	c.dc.Rotate(angle)
}

// Scale 缩放画布
func (c *Canvas) Scale(x, y float64) {
	c.dc.Scale(x, y)
}

// Translate 平移画布
func (c *Canvas) Translate(x, y float64) {
	c.dc.Translate(x, y)
}

// Push 保存当前状态
func (c *Canvas) Push() {
	c.dc.Push()
}

// Pop 恢复之前的状态
func (c *Canvas) Pop() {
	c.dc.Pop()
}

// DrawRoundedRectangle 绘制圆角矩形
func (c *Canvas) DrawRoundedRectangle(x, y, width, height, radius float64, fill bool) {
	c.dc.DrawRoundedRectangle(x, y, width, height, radius)
	if fill {
		c.dc.Fill()
	} else {
		c.dc.Stroke()
	}
}

// ApplyMask 应用蒙版
func (c *Canvas) ApplyMask(mask image.Image) {
	// gg包只接受*image.Alpha作为蒙版
	// 这里我们需要转换mask为Alpha格式
	bounds := mask.Bounds()
	alphaMask := image.NewAlpha(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := mask.At(x, y).RGBA()
			alphaMask.SetAlpha(x, y, color.Alpha{A: uint8(a >> 8)})
		}
	}

	c.dc.SetMask(alphaMask)
}

// ClearMask 清除蒙版
func (c *Canvas) ClearMask() {
	// gg没有直接的ResetMask方法，但可以通过设置nil来清除蒙版
	c.dc.SetMask(nil)
}
