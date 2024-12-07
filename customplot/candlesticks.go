package customplot

import (
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// CandleWidthMultiplier sets the candle width relative to the DefaultLineStyle.Width.
var CandleWidthMultiplier = 3

var (
	DefaultUpColor   = color.RGBA{R: 128, G: 192, B: 128, A: 255} // Green
	DefaultDownColor = color.RGBA{R: 255, G: 128, B: 128, A: 255} // Red
	DefaultLineStyle = plotter.DefaultLineStyle
)

// CandleStick represents a candlestick plotter implementing plot.Plotter.
type CandleStick struct {
	DataProvider MarketData

	// UpColor is the color for candles where Close >= Open.
	UpColor color.Color

	// DownColor is the color for candles where Close < Open.
	DownColor color.Color

	// LineStyle defines the style of the lines used to draw candles.
	draw.LineStyle

	// CandleWidth defines the width of each candle.
	CandleWidth vg.Length

	// FixedLineColor determines if a fixed color is used for all candle lines.
	FixedLineColor bool
}

// NewCandleStick initializes a CandleStick with the provided market data.
func NewCandleStick(data MarketDataProvider) (*CandleStick, error) {
	cloneData, err := CloneMarketData(data)
	if err != nil {
		return nil, err
	}

	return &CandleStick{
		DataProvider:   cloneData,
		FixedLineColor: true,
		UpColor:        DefaultUpColor,
		DownColor:      DefaultDownColor,
		LineStyle:      DefaultLineStyle,
		CandleWidth:    vg.Length(CandleWidthMultiplier) * plotter.DefaultLineStyle.Width,
	}, nil
}

// Plot renders the candlestick chart on the given canvas and plot.
func (cs *CandleStick) Plot(cnv draw.Canvas, plt *plot.Plot) {
	transformX, transformY := plt.Transforms(&cnv)
	currentStyle := cs.LineStyle

	for _, record := range cs.DataProvider {
		var fillColor color.Color
		if record.Close >= record.Open {
			fillColor = cs.UpColor
		} else {
			fillColor = cs.DownColor
		}

		if !cs.FixedLineColor {
			currentStyle.Color = fillColor
		}

		// Convert data points to canvas coordinates.
		xPos := transformX(record.Time)
		highY := transformY(record.High)
		lowY := transformY(record.Low)
		topY := transformY(math.Max(record.Open, record.Close))
		bottomY := transformY(math.Min(record.Open, record.Close))

		// Draw the upper wick.
		upperWick := cnv.ClipLinesY([]vg.Point{{xPos, highY}, {xPos, topY}})
		cnv.StrokeLines(currentStyle, upperWick...)

		// Draw the lower wick.
		lowerWick := cnv.ClipLinesY([]vg.Point{{xPos, lowY}, {xPos, bottomY}})
		cnv.StrokeLines(currentStyle, lowerWick...)

		// Draw the candle body.
		body := cnv.ClipPolygonY([]vg.Point{
			{xPos - cs.CandleWidth/2, topY},
			{xPos + cs.CandleWidth/2, topY},
			{xPos + cs.CandleWidth/2, bottomY},
			{xPos - cs.CandleWidth/2, bottomY},
			{xPos - cs.CandleWidth/2, topY},
		})
		cnv.FillPolygon(fillColor, body)
		cnv.StrokeLines(currentStyle, body)
	}
}

// DataRange computes the data boundaries for the candlestick plot.
func (cs *CandleStick) DataRange() (minX, maxX, minY, maxY float64) {
	minX = math.Inf(1)
	maxX = math.Inf(-1)
	minY = math.Inf(1)
	maxY = math.Inf(-1)

	for _, record := range cs.DataProvider {
		if record.Time < minX {
			minX = record.Time
		}
		if record.Time > maxX {
			maxX = record.Time
		}
		if record.Low < minY {
			minY = record.Low
		}
		if record.High > maxY {
			maxY = record.High
		}
	}

	return
}

// GlyphBoxes provides the bounding boxes for glyphs in the plot.
func (cs *CandleStick) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	boxes := make([]plot.GlyphBox, 2)

	minX, maxX, minY, maxY := cs.DataRange()

	boxes[0].X = plt.X.Norm(minX)
	boxes[0].Y = plt.Y.Norm(minY)
	boxes[0].Rectangle = vg.Rectangle{
		Min: vg.Point{X: -(cs.CandleWidth + cs.LineStyle.Width) / 2, Y: 0},
		Max: vg.Point{X: 0, Y: 0},
	}

	boxes[1].X = plt.X.Norm(maxX)
	boxes[1].Y = plt.Y.Norm(maxY)
	boxes[1].Rectangle = vg.Rectangle{
		Min: vg.Point{X: 0, Y: 0},
		Max: vg.Point{X: (cs.CandleWidth + cs.LineStyle.Width) / 2, Y: 0},
	}

	return boxes
}
