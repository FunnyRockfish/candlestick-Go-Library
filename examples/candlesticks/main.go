package main

import (
	"log"

	"gonum.org/v1/plot"
	"testPlotextLib/library/custplotter"
	"testPlotextLib/library/examples"
)

func main() {
	n := 60
	fakeTOHLCVs := examples.CreateTOHLCVExampleData(n)

	p := plot.New()

	p.Title.Text = "Candlesticks"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Price"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04:05"}

	bars, err := custplotter.NewCandlesticks(fakeTOHLCVs)
	if err != nil {
		log.Panic(err)
	}

	p.Add(bars)

	err = p.Save(450, 200, "candlesticks.png")
	if err != nil {
		log.Panic(err)
	}
}