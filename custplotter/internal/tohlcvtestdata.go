package internal

import (
	"testPlotextLib/library/custplotter"
	"testPlotextLib/library/examples"
)

func CreateTOHLCVTestData() custplotter.TOHLCVs {
	return examples.CreateTOHLCVExampleData(250)
}
