package main

import (
	"log"

	"candlestick-Go-Library/customplot"
	"candlestick-Go-Library/sampledata"
	"gonum.org/v1/plot"
)

func main() {
	numEntries := 60
	mockVolumeData := sampledata.GenerateCandlestickData(numEntries)

	volumePlot := plot.New()

	volumePlot.Title.Text = "Trading Volume"
	volumePlot.X.Label.Text = "Timestamp"
	volumePlot.Y.Label.Text = "Shares"
	volumePlot.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}

	barPlot, err := customplot.InitializeVolumeBars(mockVolumeData)
	if err != nil {
		log.Fatalf("Error creating volume bars: %v", err)
	}

	volumePlot.Add(barPlot)

	err = volumePlot.Save(500, 200, "volume_graph.png")
	if err != nil {
		log.Fatalf("Error saving volume graph: %v", err)
	}
}
