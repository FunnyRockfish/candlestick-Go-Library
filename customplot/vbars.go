package customplot

import (
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// VolumeBars реализует интерфейс Plotter, создавая
// диаграмму объёма торгов.
type VolumeBars struct {
	MarketData MarketData

	// PositiveColor — цвет столбцов, где C >= O.
	PositiveColor color.Color

	// NegativeColor — цвет столбцов, где C < O.
	NegativeColor color.Color

	// LineStyle — стиль линий для отображения столбцов.
	draw.LineStyle

	// TickSize — размер отметок для открытия и закрытия.
	TickSize vg.Length
}

// InitializeVolumeBars создаёт новый экземпляр VolumeBars на основе предоставленных данных.
func InitializeVolumeBars(data MarketDataProvider) (*VolumeBars, error) {
	clonedData, err := CloneMarketData(data)
	if err != nil {
		return nil, err
	}

	// Определяем минимальный интервал времени между записями.
	minInterval := 0.0
	if len(clonedData) > 1 {
		minInterval = math.MaxFloat64
		for i := 0; i < len(clonedData)-1; i++ {
			interval := clonedData[i+1].Time - clonedData[i].Time
			if interval < minInterval {
				minInterval = interval
			}
		}
	}

	return &VolumeBars{
		MarketData:    clonedData,
		PositiveColor: color.RGBA{R: 0, G: 128, B: 0, A: 255}, // Ярко-зелёный для роста.
		NegativeColor: color.RGBA{R: 196, G: 0, B: 0, A: 255}, // Ярко-красный для падения.
		LineStyle:     plotter.DefaultLineStyle,
		TickSize:      DefaultTickSize,
	}, nil
}

// Plot отображает диаграмму объёма на заданном холсте и графике.
func (vb *VolumeBars) Plot(canvas draw.Canvas, plt *plot.Plot) {
	transformX, transformY := plt.Transforms(&canvas)
	currentStyle := vb.LineStyle

	for _, entry := range vb.MarketData {
		// Устанавливаем цвет столбца в зависимости от движения цены.
		if entry.Close >= entry.Open {
			currentStyle.Color = vb.PositiveColor
		} else {
			currentStyle.Color = vb.NegativeColor
		}

		// Преобразуем данные в координаты холста.
		xPos := transformX(entry.Time)
		baseY := transformY(0)
		volumeY := transformY(entry.Volume)

		// Рисуем вертикальный столбец объёма.
		bar := canvas.ClipLinesY([]vg.Point{{xPos, baseY}, {xPos, volumeY}})
		canvas.StrokeLines(currentStyle, bar...)
	}
}

// DataRange вычисляет диапазоны данных для масштабирования графика.
func (vb *VolumeBars) DataRange() (minX, maxX, minY, maxY float64) {
	minX = math.Inf(1)
	maxX = math.Inf(-1)
	minY = 0
	maxY = math.Inf(-1)

	for _, entry := range vb.MarketData {
		if entry.Time < minX {
			minX = entry.Time
		}
		if entry.Time > maxX {
			maxX = entry.Time
		}
		if entry.Volume > maxY {
			maxY = entry.Volume
		}
	}

	return
}

// GlyphBoxes предоставляет ограничивающие рамки для глифов на графике.
func (vb *VolumeBars) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	boxes := make([]plot.GlyphBox, 2)

	minX, maxX, minY, maxY := vb.DataRange()

	// Определяем рамку для минимальных значений.
	boxes[0].X = plt.X.Norm(minX)
	boxes[0].Y = plt.Y.Norm(minY)
	boxes[0].Rectangle = vg.Rectangle{
		Min: vg.Point{X: -vb.LineStyle.Width / 2, Y: 0},
		Max: vg.Point{X: 0, Y: 0},
	}

	// Определяем рамку для максимальных значений.
	boxes[1].X = plt.X.Norm(maxX)
	boxes[1].Y = plt.Y.Norm(maxY)
	boxes[1].Rectangle = vg.Rectangle{
		Min: vg.Point{X: 0, Y: 0},
		Max: vg.Point{X: +vb.LineStyle.Width / 2, Y: 0},
	}

	return boxes
}
