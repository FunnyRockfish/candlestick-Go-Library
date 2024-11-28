package library

import (
	"math"

	"gonum.org/v1/plot"
)

// UniteAxisRanges sets the range of all axes to the minimum and the maximum of all axes.
func UniteAxisRanges(axes []*plot.Axis) {
	minRange := math.MaxFloat64
	maxRange := -math.MaxFloat64

	for _, axis := range axes {
		minRange = math.Min(axis.Min, minRange)
		maxRange = math.Max(axis.Max, maxRange)
	}

	for _, axis := range axes {
		axis.Min = minRange
		axis.Max = maxRange
	}

	return
}
