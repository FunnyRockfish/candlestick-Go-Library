package library

import (
	"math"

	"candlestick-Go-Library/logger"
	"gonum.org/v1/plot"
)

// AlignAxisRanges приводит диапазоны всех осей к общему минимальному и максимальному значению.
func AlignAxisRanges(axes []*plot.Axis) {
	log := logger.CreateLogger() // Создание логгера

	if len(axes) == 0 {
		log.Warn("Передан пустой список осей, выравнивание не выполнено")
		return
	}

	commonMin := math.MaxFloat64
	commonMax := -math.MaxFloat64

	// Находим общий минимальный и максимальный диапазон.
	for _, axis := range axes {
		commonMin = math.Min(axis.Min, commonMin)
		commonMax = math.Max(axis.Max, commonMax)
	}

	// Устанавливаем общий диапазон для всех осей.
	for _, axis := range axes {
		axis.Min = commonMin
		axis.Max = commonMax
		log.Infof("Диапазон оси изменён: Min=%.2f, Max=%.2f", axis.Min, axis.Max)
	}

	log.Info("Выравнивание диапазонов осей выполнено успешно")
}
