package main

import (
	"os"
	"time"

	candlestickLib "candlestick-Go-Library"
	customPlotLib "candlestick-Go-Library/customplot"
	"candlestick-Go-Library/logger"
	sampleDataLib "candlestick-Go-Library/sampledata"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type TimeFormatter struct {
	Layout   string
	Interval time.Duration
}

func (tf TimeFormatter) Ticks(minVal, maxVal float64) []plot.Tick {
	startTime := time.Unix(int64(minVal), 0)
	endTime := time.Unix(int64(maxVal), 0)

	var tickMarks []plot.Tick
	for current := startTime; current.Before(endTime); current = current.Add(tf.Interval) {
		tickMarks = append(tickMarks, plot.Tick{
			Value: float64(current.Unix()),
			Label: current.Format(tf.Layout),
		})
	}
	return tickMarks
}

func main() {
	log := logger.CreateLogger()
	dataPoints := 50
	ohlcvData := sampleDataLib.GenerateCandlestickData(dataPoints)

	mainPlot := plot.New()
	mainPlot.Title.Text = "Канделябры и Объёмные Барры"
	mainPlot.Y.Label.Text = "Цена"
	mainPlot.X.Tick.Marker = TimeFormatter{
		Layout:   "2006-01-02\n15:04:05",
		Interval: time.Hour, // например: метка каждый час
	}

	candles, err := customPlotLib.NewCandleStick(ohlcvData)
	if err != nil {
		log.Error("Ошибка при создании свечей: ", err)
	}

	mainPlot.Add(candles)

	volumePlot := plot.New()
	volumePlot.X.Label.Text = "Время"
	volumePlot.Y.Label.Text = "Объём"
	volumePlot.X.Tick.Marker = TimeFormatter{
		Layout:   "2006-01-02\n15:04:05",
		Interval: time.Hour, // например: метка каждый час
	}

	volumes, err := customPlotLib.InitializeVolumeBars(ohlcvData)
	if err != nil {
		log.Error("Ошибка при инициализации объёмов: ", err)
	}

	volumePlot.Add(volumes)

	candlestickLib.AlignAxisRanges([]*plot.Axis{&mainPlot.X, &volumePlot.X})

	layout := candlestickLib.GridLayout{
		RowSizes:    []float64{2, 1}, // 2/3
		ColumnSizes: []float64{1},
	}

	plotCollection := [][]*plot.Plot{{mainPlot}, {volumePlot}}

	imageCanvas := vgimg.New(2250, 300)
	drawingContext := draw.New(imageCanvas)

	alignedCanvases := layout.AlignCanvases(plotCollection, drawingContext)
	plotCollection[0][0].Draw(alignedCanvases[0][0])
	plotCollection[1][0].Draw(alignedCanvases[1][0])

	outputFile := "output_chart.png"
	fileHandle, err := os.Create(outputFile)
	if err != nil {
		log.Error("Не удалось создать файл: ", err)
	}
	defer fileHandle.Close()

	pngCanvas := vgimg.PngCanvas{Canvas: imageCanvas}
	if _, err = pngCanvas.WriteTo(fileHandle); err != nil {
		log.Error("Ошибка при записи изображения: ", err)
	}
}
