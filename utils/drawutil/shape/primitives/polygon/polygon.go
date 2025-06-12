package polygon

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/jiushengTech/common/utils/drawutil/colorx"
	"github.com/jiushengTech/common/utils/drawutil/shape/base"
)

// Polygon 表示多边形图形
type Polygon struct {
	base.BaseShape
	Fill bool `json:"fill"` // 是否填充
}

// New 创建一个新的多边形
func New(fill bool, options ...base.Option) *Polygon {
	polygon := &Polygon{
		BaseShape: base.BaseShape{
			ShapeType: "polygon",
			Points:    []*base.Point{},
			Color:     colorx.Black,
			LineWidth: 2.0,
		},
		Fill: fill,
	}

	// 应用所有选项
	base.ApplyOptions(&polygon.BaseShape, options...)

	return polygon
}

// Draw 实现Shape接口的绘制方法
func (p *Polygon) Draw(dc *gg.Context, width, height float64) error {
	if len(p.Points) < 3 {
		return fmt.Errorf("多边形至少需要3个点")
	}

	// 设置颜色和线宽
	dc.SetColor(p.Color)
	dc.SetLineWidth(p.LineWidth)

	// 开始绘制路径
	dc.MoveTo(p.Points[0].X, p.Points[0].Y)

	// 添加所有点到路径
	for i := 1; i < len(p.Points); i++ {
		dc.LineTo(p.Points[i].X, p.Points[i].Y)
	}

	// 闭合路径
	dc.ClosePath()

	// 绘制
	if p.Fill {
		dc.Fill()
	} else {
		dc.Stroke()
	}

	return nil
}
