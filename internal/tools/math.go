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
	result := AddVectors(vectors...)
	return MultiplyVectorByNum(1.0/float64(len(vectors)), result)
}
func MultiplyVectorByNum(f float64, v []float64) []float64 {
	for i := range v {
		v[i] *= f
	}
	return v
}
func AddVectors(vec ...[]float64) []float64 {
	max_len := 0
	for _, v := range vec {
		if len(v) > max_len {
			max_len = len(v)
		}
	}

	result_vec := make([]float64, max_len)
	for _, v := range vec {
		for i, f := range v {
			result_vec[i] += f
		}
	}

	return result_vec
}

type Matrix struct {
	Dict map[string]int
	Vec  []float64
	Len  int
}

func MergeMatrix(a, b Matrix) Matrix {
	for k := range b.Dict {
		if a.Dict[k] == 0 {
			a.Dict[k] = len(a.Dict) + 1
		}
	}

	//Merge Vec of Ingredient
	a.Vec = append(a.Vec, make([]float64, len(a.Dict)-len(a.Vec))...)
	b_vec := make([]float64, len(a.Dict))
	for k, v := range b.Dict {
		i := a.Dict[k]
		b_vec[i-1] = b.Vec[v-1]
	}

	vec := AddVectors(MultiplyVectorByNum(float64(a.Len), a.Vec), MultiplyVectorByNum(float64(b.Len), b_vec))
	a.Vec = MultiplyVectorByNum(1.0/float64(a.Len+1), vec)

	return a
}
