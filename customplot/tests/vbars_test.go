package tests_test

import (
	"testing"

	"candlestick-Go-Library/customplot"
	"candlestick-Go-Library/customplot/internal"
	"candlestick-Go-Library/logger"
	"gonum.org/v1/plot"
)

func TestNewVBars(t *testing.T) {
	log := logger.CreateLogger()

	log.Info("Generating market data for vertical bars test")
	testData := internal.GenerateMarkerData()

	log.Info("Creating new plot for vertical bars")
	plotInstance := plot.New()
	plotInstance.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04:05"}

	log.Info("Creating vertical bars")
	bars, err := customplot.InitializeVolumeBars(testData)
	if err != nil {
		log.Errorf("Failed to create vertical bars: %v", err)
		t.FailNow()
	}

	plotInstance.Add(bars)

	testFile := "testdata/volume_graph.png"
	log.Infof("Saving vertical bars plot to file: %s", testFile)
	err = plotInstance.Save(1180, 200, testFile)
	if err != nil {
		log.Errorf("Failed to save vertical bars plot: %v", err)
		t.FailNow()
	}

	log.Info("Validating generated vertical bars image")
	internal.TestImage(t, testFile)
	log.Info("Vertical bars test passed successfully")
}
