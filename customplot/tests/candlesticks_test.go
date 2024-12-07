package tests_test

import (
	"testing"

	"candlestick-Go-Library/customplot"
	"candlestick-Go-Library/customplot/internal"
	"candlestick-Go-Library/logger"
	"gonum.org/v1/plot"
)

func TestNewCandlesticks(t *testing.T) {
	log := logger.CreateLogger()

	log.Info("Generating market data for test")
	testData := internal.GenerateMarkerData()

	log.Info("Creating new plot")
	plotInstance := plot.New()
	plotInstance.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04:05"}

	log.Info("Creating candlestick bars")
	bars, err := customplot.NewCandleStick(testData)
	if err != nil {
		log.Errorf("Failed to create candlestick bars: %v", err)
		t.FailNow()
	}
	plotInstance.Add(bars)

	testFile := "testdata/candlestick_chart.png"
	log.Infof("Saving plot to file: %s", testFile)
	err = plotInstance.Save(1000, 150, testFile)
	if err != nil {
		log.Errorf("Failed to save plot: %v", err)
		t.FailNow()
	}

	log.Info("Validating generated image")
	internal.TestImage(t, testFile)
	log.Info("Test passed successfully")
}
