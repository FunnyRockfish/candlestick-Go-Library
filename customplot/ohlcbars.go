package customplot

import (
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// DefaultTickSize defines the length of the open and close ticks.
var DefaultTickSize = vg.Points(2)

// PriceBars implements the plot.Plotter interface, generating
// a bar chart from time, open, high, low, close data points.
type PriceBars struct {
	MarketData MarketData

	// PositiveColor is the color for bars where Close >= Open.
	PositiveColor color.Color

	// NegativeColor is the color for bars where Close < Open.
	NegativeColor color.Color

	// LineStyle specifies the styling for the bar lines.
	draw.LineStyle

	// TickSize determines the size of ticks at the open and close points.
	TickSize vg.Length
}

// InitializePriceBars creates a new PriceBars instance using the provided market data.
func InitializePriceBars(data MarketDataProvider) (*PriceBars, error) {
	clonedData, err := CloneMarketData(data)
	if err != nil {
		return nil, err
	}

	return &PriceBars{
		MarketData:    clonedData,
		PositiveColor: color.RGBA{R: 0, G: 128, B: 0, A: 255}, // Green for positive movement.
		NegativeColor: color.RGBA{R: 196, G: 0, B: 0, A: 255}, // Red for negative movement.
		LineStyle:     plotter.DefaultLineStyle,
		TickSize:      DefaultTickSize,
	}, nil
}

// Plot renders the price bars on the given canvas and plot.
func (pb *PriceBars) Plot(cnv draw.Canvas, plt *plot.Plot) {
	transformX, transformY := plt.Transforms(&cnv)
	currentStyle := pb.LineStyle

	for _, record := range pb.MarketData {
		// Set the line color based on price movement.
		if record.Close >= record.Open {
			currentStyle.Color = pb.PositiveColor
		} else {
			currentStyle.Color = pb.NegativeColor
		}

		// Convert data points to canvas coordinates.
		xPos := transformX(record.Time)
		openY := transformY(record.Open)
		highY := transformY(record.High)
		lowY := transformY(record.Low)
		closeY := transformY(record.Close)

		// Draw the high-low vertical line.
		verticalLine := cnv.ClipLinesY([]vg.Point{{xPos, lowY}, {xPos, highY}})
		cnv.StrokeLines(currentStyle, verticalLine...)

		// Draw the open tick.
		if cnv.Contains(vg.Point{X: xPos, Y: openY}) {
			cnv.StrokeLine2(currentStyle, xPos, openY, xPos-pb.TickSize, openY)
		}

		// Draw the close tick.
		if cnv.Contains(vg.Point{X: xPos, Y: closeY}) {
			cnv.StrokeLine2(currentStyle, xPos, closeY, xPos+pb.TickSize, closeY)
		}
	}
}

// DataRange calculates the boundaries of the data for scaling the plot.
func (pb *PriceBars) DataRange() (minX, maxX, minY, maxY float64) {
	minX = math.Inf(1)
	maxX = math.Inf(-1)
	minY = math.Inf(1)
	maxY = math.Inf(-1)

	for _, record := range pb.MarketData {
		minX = math.Min(minX, record.Time)
		maxX = math.Max(maxX, record.Time)
		minY = math.Min(minY, record.Low)
		maxY = math.Max(maxY, record.High)
	}

	return
}

// GlyphBoxes provides bounding boxes for glyphs in the plot.
func (pb *PriceBars) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	boxes := make([]plot.GlyphBox, 2)

	minX, maxX, minY, maxY := pb.DataRange()

	boxes[0].X = plt.X.Norm(minX)
	boxes[0].Y = plt.Y.Norm(minY)
	boxes[0].Rectangle = vg.Rectangle{
		Min: vg.Point{X: -pb.TickSize},
		Max: vg.Point{},
	}

	boxes[1].X = plt.X.Norm(maxX)
	boxes[1].Y = plt.Y.Norm(maxY)
	boxes[1].Rectangle = vg.Rectangle{
		Min: vg.Point{},
		Max: vg.Point{X: +pb.TickSize},
	}

	return boxes
}
