package test

import (
	"reflect"
	"testing"

	"github.com/madswillem/recipeApp_Backend_Go/internal/tools"
)

func TestCalculateAverage(t *testing.T) {
	t.Run("average positive numbers", func(t *testing.T) {
		avg := tools.CalculateAverage([]float64{1, 2, 3, 4, 5})
		if avg != 3 {
			t.Errorf("Expected average of 3, got %f", avg)
		}
	})
	t.Run("average of positive floats", func(t *testing.T) {
		avg := tools.CalculateAverage([]float64{1.1, 2.2, 3.3, 4.4, 5.5})
		if avg != 3.3 {
			t.Errorf("Expected average of 3.3, got %f", avg)
		}
	})
	t.Run("average of negative numbers", func(t *testing.T) {
		avg := tools.CalculateAverage([]float64{-1, -2, -3, -4, -5})
		if avg != -3 {
			t.Errorf("Expected average of -3, got %f", avg)
		}
	})
	t.Run("average of negative floats", func(t *testing.T) {
		avg := tools.CalculateAverage([]float64{-1.1, -2.2, -3.3, -4.4, -5.5})
		if avg != -3.3 {
			t.Errorf("Expected average of -3.3, got %f", avg)
		}
	})
	t.Run("average of negative and positive numbers", func(t *testing.T) {
		avg := tools.CalculateAverage([]float64{-1, 2, -3, 4, -5})
		if avg != -0.6 {
			t.Errorf("Expected average of 0, got %f", avg)
		}
	})
}

func TestGetCurrentData(t *testing.T) {
	_, err := tools.GetCurrentData()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestAverageVectors(t *testing.T) {
	v1 := []float64{1.0, 2.0, 3.0}
	v2 := []float64{4.0, 5.0}
	v3 := []float64{7.0, 8.0, 9.0, 10.0}

	expected := []float64{4, 5, 6, 3.3333333333333335}

	avg := tools.AverageVectors(v1, v2, v3)

	if reflect.DeepEqual(avg, expected) {
		t.Errorf("Expected %v but got %v", expected, avg)
	}
}
