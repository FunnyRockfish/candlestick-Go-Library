package tests_test

import (
	"testing"

	"candlestick-Go-Library/customplot"
	"candlestick-Go-Library/customplot/internal"
	"candlestick-Go-Library/logger"
	"gonum.org/v1/plot"
)

func TestNewOHLCBars(t *testing.T) {
	log := logger.CreateLogger()

	log.Info("Generating market data for OHLC bars test")
	testData := internal.GenerateMarkerData()

	log.Info("Creating new plot for OHLC bars")
	plotInstance := plot.New()
	plotInstance.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04:05"}

	log.Info("Creating OHLC bars")
	bars, err := customplot.InitializePriceBars(testData)
	if err != nil {
		log.Errorf("Failed to create OHLC bars: %v", err)
		t.FailNow()
	}

	plotInstance.Add(bars)

	testFile := "testdata/ohlc_graph.png"
	log.Infof("Saving OHLC bars plot to file: %s", testFile)
	err = plotInstance.Save(1180, 200, testFile)
	if err != nil {
		log.Errorf("Failed to save OHLC bars plot: %v", err)
		t.FailNow()
	}

	log.Info("Validating generated OHLC bars image")
	internal.TestImage(t, testFile)
	log.Info("OHLC bars test passed successfully")
}
