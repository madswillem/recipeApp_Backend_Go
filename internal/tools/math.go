package tools

import (
	"math"
)

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
func PercentageCalculator(rating float64, percentage float64) float64 {
	g := rating / 100.0
	h := g * percentage

	return RoundFloat(h, 1)
}
func CalculateAverage(numbers []float64) float64 {
	sum := 0.0
	count := len(numbers)

	if count == 0 {
		return 0.0
	}

	for _, num := range numbers {
		sum += num
	}

	return sum / float64(count)
}

func AverageVectors(vectors ...[]float64) []float64 {
	if len(vectors) == 0 {
		return nil
	}

	// Find the length of the longest vector
	maxLength := 0
	for _, v := range vectors {
		if len(v) > maxLength {
			maxLength = len(v)
		}
	}

	avg := make([]float64, maxLength)
	for _, v := range vectors {
		for i := 0; i < len(v); i++ {
			avg[i] += v[i]
		}
	}

	numVectors := float64(len(vectors))
	for i := 0; i < maxLength; i++ {
		avg[i] /= numVectors
	}

	return avg
}
