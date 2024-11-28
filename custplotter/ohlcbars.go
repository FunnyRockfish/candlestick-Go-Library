package custplotter

import (
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// DefaultTickWidth задаёт ширину отметок open и close.
var DefaultTickWidth = vg.Points(2)

// OHLCBars реализует интерфейс Plotter, создавая
// столбчатую диаграмму из кортежей time, open, high, low, close.
type OHLCBars struct {
	MarketData

	// ColorUp — цвет столбцов, где C >= O.
	ColorUp color.Color

	// ColorDown — цвет столбцов, где C < O.
	ColorDown color.Color

	// LineStyle — стиль линий для рисования столбцов.
	draw.LineStyle

	// TickWidth — ширина отметок, рисуемых сверху
	// и снизу каждой полосы ошибки.
	TickWidth vg.Length
}

// InitializeOHLCBars создаёт новый объект для
// построения диаграммы столбцов на основе предоставленных данных.
func InitializeOHLCBars(data MarketDataProvider) (*OHLCBars, error) {
	cpy, err := CloneMarketData(data)
	if err != nil {
		return nil, err
	}

	return &OHLCBars{
		MarketData: cpy,
		ColorUp:    color.RGBA{R: 0, G: 128, B: 0, A: 255}, // Зелёный цвет более заметен для глаза.
		ColorDown:  color.RGBA{R: 196, G: 0, B: 0, A: 255},
		LineStyle:  plotter.DefaultLineStyle,
		TickWidth:  DefaultTickWidth,
	}, nil
}

// Plot реализует метод Plot интерфейса plot.Plotter.
func (bars *OHLCBars) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)
	lineStyle := bars.LineStyle

	for _, TOHLCV := range bars.MarketData {
		if TOHLCV.Close >= TOHLCV.Open {
			lineStyle.Color = bars.ColorUp
		} else {
			lineStyle.Color = bars.ColorDown
		}

		// Преобразование данных в соответствующие координаты для рисования.
		x := trX(TOHLCV.Time)
		yo := trY(TOHLCV.Open)
		yh := trY(TOHLCV.High)
		yl := trY(TOHLCV.Low)
		yc := trY(TOHLCV.Close)

		// Рисование вертикальной полосы (high-low).
		bar := c.ClipLinesY([]vg.Point{{x, yl}, {x, yh}})
		c.StrokeLines(lineStyle, bar...)

		// Рисование отметки для open.
		if c.Contains(vg.Point{X: x, Y: yo}) {
			c.StrokeLine2(lineStyle, x, yo, x-bars.TickWidth, yo)
		}

		// Рисование отметки для close.
		if c.Contains(vg.Point{X: x, Y: yc}) {
			c.StrokeLine2(lineStyle, x, yc, x+bars.TickWidth, yc)
		}

	}
}

// DataRange реализует метод DataRange интерфейса plot.DataRanger.
func (bars *OHLCBars) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = math.Inf(1)
	xmax = math.Inf(-1)
	ymin = math.Inf(1)
	ymax = math.Inf(-1)
	for _, TOHLCV := range bars.MarketData {
		xmin = math.Min(xmin, TOHLCV.Time)
		xmax = math.Max(xmax, TOHLCV.Time)
		ymin = math.Min(ymin, TOHLCV.Low)
		ymax = math.Max(ymax, TOHLCV.High)
	}
	return
}

// GlyphBoxes реализует метод GlyphBoxes интерфейса plot.GlyphBoxer.
func (bars *OHLCBars) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	boxes := make([]plot.GlyphBox, 2)

	xmin, xmax, ymin, ymax := bars.DataRange()

	boxes[0].X = plt.X.Norm(xmin)
	boxes[0].Y = plt.Y.Norm(ymin)
	boxes[0].Rectangle = vg.Rectangle{
		Min: vg.Point{X: -bars.TickWidth},
		Max: vg.Point{},
	}

	boxes[1].X = plt.X.Norm(xmax)
	boxes[1].Y = plt.Y.Norm(ymax)
	boxes[1].Rectangle = vg.Rectangle{
		Min: vg.Point{},
		Max: vg.Point{X: +bars.TickWidth},
	}

	return boxes
}
