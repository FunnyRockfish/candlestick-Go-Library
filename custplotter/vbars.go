package custplotter

import (
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// VBars реализует интерфейс Plotter, создавая
// столбчатую диаграмму объёма.
type VBars struct {
	MarketData

	// ColorUp — цвет столбцов, где C >= O.
	ColorUp color.Color

	// ColorDown — цвет столбцов, где C < O.
	ColorDown color.Color

	// LineStyle — стиль линий для рисования столбцов.
	draw.LineStyle
}

// InitializeVBars создаёт новый объект для построения
// столбчатой диаграммы на основе предоставленных данных.
func InitializeVBars(data MarketDataProvider) (*VBars, error) {
	cpy, err := CloneMarketData(data)
	if err != nil {
		return nil, err
	}

	// Определяем минимальное значение DeltaT.
	minDeltaT := 0.0
	if len(cpy) > 1 {
		minDeltaT = math.MaxFloat64
		for i := 0; i < len(cpy)-1; i++ {
			minDeltaT = math.Min(minDeltaT, cpy[i+1].Time-cpy[i].Time)
		}
	}

	return &VBars{
		MarketData: cpy,
		ColorUp:    color.RGBA{R: 0, G: 128, B: 0, A: 255}, // Зелёный цвет более заметен для глаза.
		ColorDown:  color.RGBA{R: 196, G: 0, B: 0, A: 255},
		LineStyle:  plotter.DefaultLineStyle,
	}, nil
}

// Plot реализует метод Plot интерфейса plot.Plotter.
func (bars *VBars) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)
	lineStyle := bars.LineStyle

	for _, TOHLCV := range bars.MarketData {
		if TOHLCV.Close >= TOHLCV.Open {
			lineStyle.Color = bars.ColorUp
		} else {
			lineStyle.Color = bars.ColorDown
		}

		// Преобразование данных в соответствующие координаты на холсте.
		x := trX(TOHLCV.Time)
		y0 := trY(0)
		y := trY(TOHLCV.Volume)

		// Рисуем столбец объёма.
		bar := c.ClipLinesY([]vg.Point{{x, y0}, {x, y}})
		c.StrokeLines(lineStyle, bar...)
	}
}

// DataRange реализует метод DataRange интерфейса plot.DataRanger.
func (bars *VBars) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = math.Inf(1)
	xmax = math.Inf(-1)
	ymin = 0
	ymax = math.Inf(-1)
	for _, TOHLCV := range bars.MarketData {
		xmin = math.Min(xmin, TOHLCV.Time)
		xmax = math.Max(xmax, TOHLCV.Time)
		ymax = math.Max(ymax, TOHLCV.Volume)
	}
	return
}

// GlyphBoxes реализует метод GlyphBoxes интерфейса plot.GlyphBoxer.
func (bars *VBars) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	boxes := make([]plot.GlyphBox, 2)

	xmin, xmax, ymin, ymax := bars.DataRange()

	// Определяем рамки для минимальных значений.
	boxes[0].X = plt.X.Norm(xmin)
	boxes[0].Y = plt.Y.Norm(ymin)
	boxes[0].Rectangle = vg.Rectangle{
		Min: vg.Point{X: -bars.LineStyle.Width / 2, Y: 0},
		Max: vg.Point{X: 0, Y: 0},
	}

	// Определяем рамки для максимальных значений.
	boxes[1].X = plt.X.Norm(xmax)
	boxes[1].Y = plt.Y.Norm(ymax)
	boxes[1].Rectangle = vg.Rectangle{
		Min: vg.Point{X: 0, Y: 0},
		Max: vg.Point{X: +bars.LineStyle.Width / 2, Y: 0},
	}

	return boxes
}
