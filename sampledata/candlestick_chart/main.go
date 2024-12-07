package main

import (
	"log"

	"candlestick-Go-Library/customplot"
	"candlestick-Go-Library/sampledata"
	"gonum.org/v1/plot"
)

func main() {
	entryCount := 10

	mockCandlestickData := sampledata.GenerateCandlestickData(entryCount)

	candlePlot := plot.New()

	candlePlot.Title.Text = "Candlestick Chart"
	candlePlot.X.Label.Text = "Datetime"
	candlePlot.Y.Label.Text = "Price"
	candlePlot.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}

	candles, err := customplot.NewCandleStick(mockCandlestickData)
	if err != nil {
		log.Fatalf("Failed to build candlestick series: %v", err)
	}

	candlePlot.Add(candles)

	err = candlePlot.Save(500, 250, "candlestick_chart.png")
	if err != nil {
		log.Fatalf("Failed to save candlestick chart: %v", err)
	}
}
