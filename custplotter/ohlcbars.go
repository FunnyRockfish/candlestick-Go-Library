package custplotter

import (
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// DefaultTickWidth is the default width of the open and close ticks.
var DefaultTickWidth = vg.Points(2)

// OHLCBars implements the Plotter interface, drawing
// a bar plot of time, open, high, low, close tuples.
type OHLCBars struct {
	TOHLCVs

	// ColorUp is the color of bars where C >= O
	ColorUp color.Color

	// ColorDown is the color of bars where C < O
	ColorDown color.Color

	// LineStyle is the style used to draw the bars.
	draw.LineStyle

	// CapWidth is the width of the caps drawn at the top
	// of each error bar.
	TickWidth vg.Length
}

// NewOHLCBars NewBars creates as new bar plotter for
// the given data.
func NewOHLCBars(TOHLCV TOHLCVer) (*OHLCBars, error) {
	cpy, err := CopyTOHLCVs(TOHLCV)
	if err != nil {
		return nil, err
	}

	return &OHLCBars{
		TOHLCVs:   cpy,
		ColorUp:   color.RGBA{R: 0, G: 128, B: 0, A: 255}, // eye is more sensible to green
		ColorDown: color.RGBA{R: 196, G: 0, B: 0, A: 255},
		LineStyle: plotter.DefaultLineStyle,
		TickWidth: DefaultTickWidth,
	}, nil
}

// Plot implements the Plot method of the plot.Plotter interface.
func (bars *OHLCBars) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)
	lineStyle := bars.LineStyle

	for _, TOHLCV := range bars.TOHLCVs {
		if TOHLCV.C >= TOHLCV.O {
			lineStyle.Color = bars.ColorUp
		} else {
			lineStyle.Color = bars.ColorDown
		}

		// Transform the data
		// to the corresponding drawing coordinate.
		x := trX(TOHLCV.T)
		yo := trY(TOHLCV.O)
		yh := trY(TOHLCV.H)
		yl := trY(TOHLCV.L)
		yc := trY(TOHLCV.C)

		bar := c.ClipLinesY([]vg.Point{{x, yl}, {x, yh}})
		c.StrokeLines(lineStyle, bar...)

		if c.Contains(vg.Point{X: x, Y: yo}) {
			c.StrokeLine2(lineStyle, x, yo, x-bars.TickWidth, yo)
		}

		if c.Contains(vg.Point{X: x, Y: yc}) {
			c.StrokeLine2(lineStyle, x, yc, x+bars.TickWidth, yc)
		}

	}
}

// DataRange implements the DataRange method
// of the plot.DataRanger interface.
func (bars *OHLCBars) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = math.Inf(1)
	xmax = math.Inf(-1)
	ymin = math.Inf(1)
	ymax = math.Inf(-1)
	for _, TOHLCV := range bars.TOHLCVs {
		xmin = math.Min(xmin, TOHLCV.T)
		xmax = math.Max(xmax, TOHLCV.T)
		ymin = math.Min(ymin, TOHLCV.L)
		ymax = math.Max(ymax, TOHLCV.H)
	}
	return
}

// GlyphBoxes implements the GlyphBoxes method
// of the plot.GlyphBoxer interface.
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
