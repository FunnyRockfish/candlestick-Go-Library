package internal

import (
	"candlestick-Go-Library/customplot"
	"candlestick-Go-Library/sampledata"
)

func GenerateMarkerData() customplot.MarketDataProvider {
	return sampledata.GenerateCandlestickData(250)
}
