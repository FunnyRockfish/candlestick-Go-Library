package internal

import (
	"candlestick-Go-Library/custplotter"
	"candlestick-Go-Library/examples"
)

func GenerateMarkerData() custplotter.MarketDataProvider {
	return examples.CreateTOHLCVExampleData(250)
}
