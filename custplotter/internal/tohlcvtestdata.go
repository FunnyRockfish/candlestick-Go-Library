package internal

import (
	"candlestick-Go-Library/custplotter"
	"candlestick-Go-Library/examples"
)

func CreateTOHLCVTestData() custplotter.TOHLCVs {
	return examples.CreateTOHLCVExampleData(250)
}
