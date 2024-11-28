package examples

import (
	"math"
	"math/rand"
	"time"

	"candlestick-Go-Library/custplotter"
)

// CreateTOHLCVExampleData генерирует и возвращает искусственные данные TOHLCV для тестирования и демонстрации.
func CreateTOHLCVExampleData(n int) custplotter.MarketData {
	rnd := rand.New(rand.NewSource(1))
	m := 4 * n
	fract := make([]float64, m)
	for i := 0; i < m; i++ {
		fract[i] = 100
	}
	stat1 := 0.0
	stat2 := 0.0
	for k := m; k > 0; k = k / 2 {
		j := 0
		for i := 0; i < m; i++ {
			if j == 0 {
				j = k
				stat2 = stat1
				stat1 = 10.0 * (float64(k)/float64(m) + 0.02) * (2.0*rnd.Float64() - 1.0)
			}
			fract[i] += float64(k-j)/float64(k)*stat1 + float64(j)/float64(k)*stat2
			j--
		}
	}

	data := make(custplotter.MarketData, n)

	loc, _ := time.LoadLocation("America/New_York")
	for i := range data {
		data[i].Time = float64(time.Date(2024, 10, 30, 03, 04, 05, 0, loc).Add(time.Duration(i) * time.Minute).Unix())
		data[i].Open = fract[4*i]
		data[i].High = math.Max(math.Max(fract[4*i], fract[4*i+1]), math.Max(fract[4*i+2], fract[4*i+3]))
		data[i].Low = math.Min(math.Min(fract[4*i], fract[4*i+1]), math.Min(fract[4*i+2], fract[4*i+3]))
		data[i].Close = fract[4*i+3]

		data[i].Volume = (data[i].High - data[i].Low + math.Abs(data[i].Close-data[i].Open)) * 100
	}
	return data
}
