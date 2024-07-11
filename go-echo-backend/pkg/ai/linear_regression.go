package ai

import (
	"gonum.org/v1/gonum/stat"
	"math"
)

func PredictSeries(series []float64, next int) (nextSeries []float64) {
	x := make([]float64, len(series))
	y := make([]float64, len(series))
	for i, num := range series {
		x[i] = float64(i)
		y[i] = num
	}
	alpha, beta := stat.LinearRegression(x, y, nil, false)

	for i := 0; i < next; i++ {
		xValue := float64(len(series) + i)
		yValue := alpha + beta*xValue
		yValue = math.Floor(yValue*100) / 100
		series = append(series, yValue)
	}
	return curvy(series[len(series)-next:])
}

func curvy(series []float64) []float64 {
	for i := 0; i < len(series); i++ {
		if i%2 == 0 {
			series[i] = 1.2 * series[i]
		} else {
			series[i] = 0.8 * series[i]
		}
	}
	return series
}
