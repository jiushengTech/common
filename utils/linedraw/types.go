package linedraw

// Point 表示二维坐标点
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// LineType 线条类型
type LineType string

// 支持的线条类型
const (
	VerticalLine   LineType = "vertical"   // 竖线
	HorizontalLine LineType = "horizontal" // 横线
)

// Line 表示一条线及其相关属性
type Line struct {
	Type      LineType   `json:"type"`       // 线条类型
	Points    []Point    `json:"points"`     // 点集合（至少需要2个点）
	Values    []float64  `json:"values"`     // 点之间的值（长度比点少1）
	Color     [3]float64 `json:"color"`      // 线条颜色，RGB值(0-1)
	LineWidth float64    `json:"line_width"` // 线条宽度
}

// 颜色常量
var (
	ColorWhite  = [3]float64{1, 1, 1} // 白色
	ColorBlack  = [3]float64{0, 0, 0} // 黑色
	ColorRed    = [3]float64{1, 0, 0} // 红色
	ColorBlue   = [3]float64{0, 0, 1} // 蓝色
	ColorGreen  = [3]float64{0, 1, 0} // 绿色
	ColorYellow = [3]float64{1, 1, 0} // 黄色
)

// 默认输出文件名
const DefaultOutputName = "result.png"
