package tools

import "math"

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
