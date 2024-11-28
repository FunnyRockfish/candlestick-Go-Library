package library

import (
	"candlestick-Go-Library/logger"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// GridLayout создаёт таблицу с субхолстами из Canvas.
// В отличие от tiles из gonum.org/v1/plot, строки и столбцы таблицы
// могут иметь разную высоту и ширину соответственно.
type GridLayout struct {
	RowSizes []float64
	// ColumnSizes задаёт количество столбцов и их относительные размеры.
	ColumnSizes []float64
	// PaddingTop, PaddingBottom, PaddingRight и PaddingLeft задают отступы
	// с соответствующей стороны таблицы.
	PaddingTop, PaddingBottom, PaddingRight, PaddingLeft vg.Length
	// PaddingX и PaddingY задают отступы между столбцами и строками соответственно.
	PaddingX, PaddingY vg.Length
}

// SubCanvas возвращает субхолст внутри canvas, который соответствует
// ячейке в столбце colIndex и строке rowIndex, где 0, 0 — верхний правый угол.
func (grid GridLayout) SubCanvas(canvas draw.Canvas, colIndex, rowIndex int) draw.Canvas {
	log := logger.CreateLogger() // Предполагается, что logger.CreateLogger() возвращает SugaredLogger

	// Проверяем корректность индексов.
	if colIndex >= len(grid.ColumnSizes) || rowIndex >= len(grid.RowSizes) {
		log.Errorf("Индекс ячейки вне допустимого диапазона: colIndex=%d, rowIndex=%d", colIndex, rowIndex)
		return draw.Canvas{}
	}

	// Сумма всех относительных размеров столбцов.
	totalColumnWidth := 0.0
	for _, colWidth := range grid.ColumnSizes {
		totalColumnWidth += colWidth
	}

	// Сумма относительных размеров для столбцов слева от текущего.
	columnOffset := 0.0
	for i := 0; i < colIndex; i++ {
		columnOffset += grid.ColumnSizes[i]
	}

	// Сумма всех относительных размеров строк.
	totalRowHeight := 0.0
	for _, rowHeight := range grid.RowSizes {
		totalRowHeight += rowHeight
	}

	// Сумма относительных размеров для строк выше текущей.
	rowOffset := 0.0
	for i := 0; i < rowIndex; i++ {
		rowOffset += grid.RowSizes[i]
	}

	// Вычисление размеров холста.
	rowUnitHeight := (canvas.Max.Y - canvas.Min.Y - grid.PaddingTop - grid.PaddingBottom -
		vg.Length(len(grid.RowSizes)-1)*grid.PaddingY) / vg.Length(totalRowHeight)
	columnUnitWidth := (canvas.Max.X - canvas.Min.X - grid.PaddingLeft - grid.PaddingRight -
		vg.Length(len(grid.ColumnSizes)-1)*grid.PaddingX) / vg.Length(totalColumnWidth)

	yMax := canvas.Max.Y - grid.PaddingTop - vg.Length(rowIndex)*grid.PaddingY - vg.Length(rowOffset)*rowUnitHeight
	yMin := yMax - vg.Length(grid.RowSizes[rowIndex])*rowUnitHeight

	xMin := canvas.Min.X + grid.PaddingLeft + vg.Length(colIndex)*grid.PaddingX + vg.Length(columnOffset)*columnUnitWidth
	xMax := xMin + vg.Length(grid.ColumnSizes[colIndex])*columnUnitWidth

	return draw.Canvas{
		Canvas: vg.Canvas(canvas),
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: xMin, Y: yMin},
			Max: vg.Point{X: xMax, Y: yMax},
		},
	}
}

// AlignCanvases возвращает двумерный массив Canvases (по строкам),
// которые обеспечивают выравнивание всех DataCanvas графиков.
func (grid GridLayout) AlignCanvases(plots [][]*plot.Plot, parentCanvas draw.Canvas) [][]draw.Canvas {
	log := logger.CreateLogger()

	if len(plots) != len(grid.RowSizes) {
		log.Errorf("Количество строк в графиках (%d) не соответствует количеству строк в таблице (%d)", len(plots), len(grid.RowSizes))
		return nil
	}

	output := make([][]draw.Canvas, len(plots))

	for rowIndex := 0; rowIndex < len(grid.RowSizes); rowIndex++ {
		if len(plots[rowIndex]) != len(grid.ColumnSizes) {
			log.Errorf("Количество столбцов в строке %d графиков (%d) не соответствует количеству столбцов в таблице (%d)", rowIndex, len(plots[rowIndex]), len(grid.ColumnSizes))
			return nil
		}

		output[rowIndex] = make([]draw.Canvas, len(plots[rowIndex]))
		for colIndex := 0; colIndex < len(grid.ColumnSizes); colIndex++ {
			output[rowIndex][colIndex] = grid.SubCanvas(parentCanvas, colIndex, rowIndex)
		}
	}

	return output
}
