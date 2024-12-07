package main

import (
	"log"
	"os"
	"time"

	"candlestick-Go-Library"
	"candlestick-Go-Library/customplot"
	"candlestick-Go-Library/sampledata"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type CustomTimeTicks struct {
	Format string
	Step   time.Duration
}

func (ctt CustomTimeTicks) Ticks(min, max float64) []plot.Tick {
	tMin := time.Unix(int64(min), 0)
	tMax := time.Unix(int64(max), 0)

	var ticks []plot.Tick
	for t := tMin; t.Before(tMax); t = t.Add(ctt.Step) {
		ticks = append(ticks, plot.Tick{
			Value: float64(t.Unix()),
			Label: t.Format(ctt.Format),
		})
	}
	return ticks
}

func main() {
	n := 300
	fakeTOHLCVs := sampledata.GenerateCandlestickData(n)

	p1 := plot.New()

	p1.Title.Text = "Candlesticks and Volume Bars"
	p1.Y.Label.Text = "Price"
	p1.X.Tick.Marker = CustomTimeTicks{
		Format: "2006-01-02\n15:04:05",
		Step:   1 * time.Hour, // например: метка каждый час
	}

	candlesticks, err := customplot.BuildCandlestickSeries(fakeTOHLCVs)
	if err != nil {
		log.Panic(err)
	}

	p1.Add(candlesticks)
	p2 := plot.New()

	p2.X.Label.Text = "Time"
	p2.Y.Label.Text = "Volume"
	p2.X.Tick.Marker = CustomTimeTicks{
		Format: "2006-01-02\n15:04:05",
		Step:   1 * time.Hour, // например: метка каждый час
	}

	vBars, err := customplot.CreateVolumeBars(fakeTOHLCVs)
	if err != nil {
		log.Panic(err)
	}

	p2.Add(vBars)

	library.AlignAxisRanges([]*plot.Axis{&p1.X, &p2.X})

	table := library.GridLayout{
		RowSizes:    []float64{2, 1}, // 2/3
		ColumnSizes: []float64{1},
	}

	plots := [][]*plot.Plot{{p1}, {p2}}

	img := vgimg.New(2250, 300)
	dc := draw.New(img)

	canvases := table.AlignCanvases(plots, dc)
	plots[0][0].Draw(canvases[0][0])
	plots[1][0].Draw(canvases[1][0])

	testFile := "align.png"
	w, err := os.Create(testFile)
	if err != nil {
		panic(err)
	}

	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(w); err != nil {
		panic(err)
	}
}