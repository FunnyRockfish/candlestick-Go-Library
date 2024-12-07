package main

import (
	"log"

	"candlestick-Go-Library/customplot"
	"candlestick-Go-Library/sampledata"
	"gonum.org/v1/plot"
)

func main() {
	dataPoints := 60
	mockData := sampledata.GenerateCandlestickData(dataPoints)

	plotInstance := plot.New()

	plotInstance.Title.Text = "OHLC Graph"
	plotInstance.X.Label.Text = "Date-Time"
	plotInstance.Y.Label.Text = "Value"
	plotInstance.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}

	ohlcPlot, err := customplot.CreateOHLCPlot(mockData)
	if err != nil {
		log.Fatalf("Failed to create OHLC plot: %v", err)
	}

	plotInstance.Add(ohlcPlot)

	err = plotInstance.Save(500, 250, "ohlc_graph.png")
	if err != nil {
		log.Fatalf("Failed to save OHLC graph: %v", err)
	}
}
