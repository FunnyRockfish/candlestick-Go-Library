package custplotter

import "gonum.org/v1/plot/plotter"

// MarketDataProvider описывает интерфейс для работы с рыночными данными.
type MarketDataProvider interface {
	// Count возвращает количество записей.
	Count() int

	// DataAt возвращает кортеж данных: время, открытие, максимум, минимум, закрытие, объем.
	DataAt(index int) (float64, float64, float64, float64, float64, float64)
}

// MarketData представляет собой срез записей с рыночными данными.
type MarketData []struct{ Time, Open, High, Low, Close, Volume float64 }

// Count реализует метод Count интерфейса MarketDataProvider.
func (md MarketData) Count() int {
	return len(md)
}

// DataAt реализует метод DataAt интерфейса MarketDataProvider.
func (md MarketData) DataAt(index int) (float64, float64, float64, float64, float64, float64) {
	return md[index].Time, md[index].Open, md[index].High, md[index].Low, md[index].Close, md[index].Volume
}

// CloneMarketData копирует данные из реализации интерфейса MarketDataProvider.
func CloneMarketData(source MarketDataProvider) (MarketData, error) {
	clone := make(MarketData, source.Count())
	for i := range clone {
		clone[i].Time, clone[i].Open, clone[i].High, clone[i].Low, clone[i].Close, clone[i].Volume = source.DataAt(i)
		if err := plotter.CheckFloats(clone[i].Open, clone[i].High, clone[i].Low, clone[i].Close, clone[i].Volume); err != nil {
			return nil, err
		}
	}
	return clone, nil
}
