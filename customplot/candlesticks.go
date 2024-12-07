package customplot

import (
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// DefaultCandleWidthFactor задаёт ширину свечи относительно DefaultLineStyle.Width.
var DefaultCandleWidthFactor = 3

var (
	DefaultColorUp   = color.RGBA{R: 128, G: 192, B: 128, A: 255} // Зеленый
	DefaultColorDown = color.RGBA{R: 255, G: 128, B: 128, A: 255} // Красный
	DefaultLineStyle = plotter.DefaultLineStyle
)

// Candlesticks реализует интерфейс Plotter, создавая
// столбчатую диаграмму из кортежей time, open, high, low, close.
type Candlesticks struct {
	MarketData MarketData

	// ColorUp — цвет свечей, где C >= O.
	ColorUp color.Color

	// ColorDown — цвет свечей, где C < O.
	ColorDown color.Color

	// LineStyle — стиль линий для рисования свечей.
	draw.LineStyle

	// CandleWidth — ширина свечи.
	CandleWidth vg.Length

	// FixedLineColor указывает, использовать ли фиксированный цвет линий для восходящих и нисходящих свечей.
	FixedLineColor bool
}

// BuildCandlestickSeries создаёт новый объект для рисования свечей
// на основе предоставленных данных.
func BuildCandlestickSeries(data MarketDataProvider) (*Candlesticks, error) {
	cpy, err := CloneMarketData(data)
	if err != nil {
		return nil, err
	}

	return &Candlesticks{
		MarketData:     cpy,
		FixedLineColor: true,
		ColorUp:        DefaultColorUp,
		ColorDown:      DefaultColorDown,
		LineStyle:      DefaultLineStyle,
		CandleWidth:    vg.Length(DefaultCandleWidthFactor) * plotter.DefaultLineStyle.Width,
	}, nil
}

// Plot реализует метод Plot интерфейса plot.Plotter.
func (sticks *Candlesticks) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)
	lineStyle := sticks.LineStyle

	for _, TOHLCV := range sticks.MarketData {
		var fillColor color.Color
		if TOHLCV.Close >= TOHLCV.Open {
			fillColor = sticks.ColorUp
		} else {
			fillColor = sticks.ColorDown
		}

		if !sticks.FixedLineColor {
			lineStyle.Color = fillColor
		}

		// Преобразование данных в соответствующие координаты на холсте.
		x := trX(TOHLCV.Time)
		yh := trY(TOHLCV.High)
		yl := trY(TOHLCV.Low)
		ymaxoc := trY(math.Max(TOHLCV.Open, TOHLCV.Close))
		yminoc := trY(math.Min(TOHLCV.Open, TOHLCV.Close))

		// Верхняя линия
		line := c.ClipLinesY([]vg.Point{{x, yh}, {x, ymaxoc}})
		c.StrokeLines(lineStyle, line...)

		// Нижняя линия
		line = c.ClipLinesY([]vg.Point{{x, yl}, {x, yminoc}})
		c.StrokeLines(lineStyle, line...)

		poly := c.ClipPolygonY([]vg.Point{{x - sticks.CandleWidth/2, ymaxoc}, {x + sticks.CandleWidth/2, ymaxoc}, {x + sticks.CandleWidth/2, yminoc}, {x - sticks.CandleWidth/2, yminoc}, {x - sticks.CandleWidth/2, ymaxoc}})
		c.FillPolygon(fillColor, poly)
		c.StrokeLines(lineStyle, poly)
	}
}

// DataRange реализует метод DataRange интерфейса plot.DataRanger.
func (sticks *Candlesticks) DataRange() (xMin, xMax, yMin, yMax float64) {
	xMin = math.Inf(1)
	xMax = math.Inf(-1)
	yMin = math.Inf(1)
	yMax = math.Inf(-1)
	for _, TOHLCV := range sticks.MarketData {
		xMin = math.Min(xMin, TOHLCV.Time)
		xMax = math.Max(xMax, TOHLCV.Time)
		yMin = math.Min(yMin, TOHLCV.Low)
		yMax = math.Max(yMax, TOHLCV.High)
	}
	return
}

// GlyphBoxes реализует метод GlyphBoxes интерфейса plot.GlyphBoxer.
func (sticks *Candlesticks) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	boxes := make([]plot.GlyphBox, 2)

	xmin, xmax, ymin, ymax := sticks.DataRange()

	boxes[0].X = plt.X.Norm(xmin)
	boxes[0].Y = plt.Y.Norm(ymin)
	boxes[0].Rectangle = vg.Rectangle{
		Min: vg.Point{X: -(sticks.CandleWidth + sticks.LineStyle.Width) / 2, Y: 0},
		Max: vg.Point{X: 0, Y: 0},
	}

	boxes[1].X = plt.X.Norm(xmax)
	boxes[1].Y = plt.Y.Norm(ymax)
	boxes[1].Rectangle = vg.Rectangle{
		Min: vg.Point{X: 0, Y: 0},
		Max: vg.Point{X: +(sticks.CandleWidth + sticks.LineStyle.Width) / 2, Y: 0},
	}

	return boxes
}
