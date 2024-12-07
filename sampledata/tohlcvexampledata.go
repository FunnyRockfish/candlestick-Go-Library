package sampledata

import (
	"math"
	"math/rand"
	"time"

	"candlestick-Go-Library/customplot"
)

// GenerateCandlestickData создает и возвращает сгенерированные данные TOHLCV для целей тестирования и демонстрации.
func GenerateCandlestickData(count int) customplot.MarketData {
	seededRand := rand.New(rand.NewSource(1))
	totalElements := 4 * count
	priceSeries := make([]float64, totalElements)

	// Инициализация базовой цены
	for idx := 0; idx < totalElements; idx++ {
		priceSeries[idx] = 100
	}

	var previousStat, currentStat float64

	// Генерация флуктуаций цен
	for step := totalElements; step > 0; step /= 2 {
		temp := 0
		for i := 0; i < totalElements; i++ {
			if temp == 0 {
				temp = step
				previousStat = currentStat
				currentStat = 10.0 * (float64(step)/float64(totalElements) + 0.02) * (2.0*seededRand.Float64() - 1.0)
			}
			priceSeries[i] += float64(step-temp)/float64(step)*currentStat + float64(temp)/float64(step)*previousStat
			temp--
		}
	}

	ohlcv := make(customplot.MarketData, count)

	nyLocation, _ := time.LoadLocation("America/New_York")
	baseTime := time.Date(2024, 10, 30, 3, 4, 5, 0, nyLocation)

	for idx := range ohlcv {
		timestamp := baseTime.Add(time.Duration(idx) * time.Minute).Unix()
		ohlcv[idx].Time = float64(timestamp)
		ohlcv[idx].Open = priceSeries[4*idx]
		ohlcv[idx].High = math.Max(math.Max(priceSeries[4*idx], priceSeries[4*idx+1]),
			math.Max(priceSeries[4*idx+2], priceSeries[4*idx+3]))
		ohlcv[idx].Low = math.Min(math.Min(priceSeries[4*idx], priceSeries[4*idx+1]),
			math.Min(priceSeries[4*idx+2], priceSeries[4*idx+3]))
		ohlcv[idx].Close = priceSeries[4*idx+3]

		// Расчет объема как функция ценовых колебаний
		ohlcv[idx].Volume = (ohlcv[idx].High - ohlcv[idx].Low + math.Abs(ohlcv[idx].Close-ohlcv[idx].Open)) * 100
	}

	return ohlcv
}
