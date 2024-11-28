package tests_test

import (
	"log"
	"testing"

	"gonum.org/v1/plot"
	"testPlotextLib/library/custplotter"
	"testPlotextLib/library/custplotter/internal"
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
	err = p.Save(180, 100, testFile)
	if err != nil {
		log.Panic(err)
	}

	internal.TestImage(t, testFile)
}
