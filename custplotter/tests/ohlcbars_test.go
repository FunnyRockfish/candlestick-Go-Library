package tests_test

import (
	"log"
	"testing"

	"candlestick-Go-Library/custplotter"
	"candlestick-Go-Library/custplotter/internal"
	"gonum.org/v1/plot"
)

func TestNewOHLCBars(t *testing.T) {
	testTOHLCVs := internal.CreateTOHLCVTestData()

	p := plot.New()

	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04:05"}

	bars, err := custplotter.NewOHLCBars(testTOHLCVs)
	if err != nil {
		log.Panic(err)
	}

	p.Add(bars)

	testFile := "testdata/ohlcbars.png"
	err = p.Save(1180, 200, testFile)
	if err != nil {
		log.Panic(err)
	}

	internal.TestImage(t, testFile)
}
