package colorx

import "image/color"

//// ---------- 常用颜色 ----------

// 颜色常量 (RGB值，范围0-1)
var (
	White   = &color.RGBA{255, 255, 255, 255}
	Black   = &color.RGBA{0, 0, 0, 255}
	Red     = &color.RGBA{255, 0, 0, 255}
	Green   = &color.RGBA{0, 255, 0, 255}
	Blue    = &color.RGBA{0, 0, 255, 255}
	Yellow  = &color.RGBA{255, 255, 0, 255}
	Cyan    = &color.RGBA{0, 255, 255, 255}
	Magenta = &color.RGBA{255, 0, 255, 255}
	Gray    = &color.RGBA{128, 128, 128, 255}
	Orange  = &color.RGBA{255, 165, 0, 255}
	Purple  = &color.RGBA{128, 0, 128, 255}
	Brown   = &color.RGBA{139, 69, 19, 255}

	// 50% 透明度版本
	Gray50  = &color.RGBA{128, 128, 128, 127}
	Blue50  = &color.RGBA{0, 0, 255, 127}
	Red50   = &color.RGBA{255, 0, 0, 127}
	Green50 = &color.RGBA{0, 255, 0, 127}
	Black50 = &color.RGBA{0, 0, 0, 127}
)
