package middleware

import (
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/models"
	"rezeptapp.ml/goApp/tools"
)

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
func percentageCalculator(rating float64, percentage float64) float64 {
	g := rating / 100.0
	h := g * percentage

	return roundFloat(h, 1)
}
func calculateAverage(numbers []float64) float64 {
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

func update(rating models.RatingStruct, change int, c *gin.Context) models.RatingStruct {

	result := rating
	data := tools.GetCurrentData(c)
	percentage := 10.0

	switch data.Day {
	case "Mon":
		result.Mon += percentageCalculator(result.Mon*float64(change), percentage)
	case "Tue":
		result.Tue += percentageCalculator(result.Tue*float64(change), percentage)
	case "Wed":
		result.Wed += percentageCalculator(result.Wed*float64(change), percentage)
	case "Thu":
		result.Thu += percentageCalculator(result.Thu*float64(change), percentage)
	case "Fri":
		result.Fri += percentageCalculator(result.Fri*float64(change), percentage)
	case "Sat":
		result.Sat += percentageCalculator(result.Sat*float64(change), percentage)
	case "Sun":
		result.Sun += percentageCalculator(result.Sun*float64(change), percentage)
	default:
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	switch data.Season {
	case "Win":
		result.Win += percentageCalculator(result.Win*float64(change), percentage)
	case "Spr":
		result.Spr += percentageCalculator(result.Spr*float64(change), percentage)
	case "Sum":
		result.Sum += percentageCalculator(result.Sum*float64(change), percentage)
	case "Aut":
		result.Aut += percentageCalculator(result.Aut*float64(change), percentage)
	default:
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	switch data.Temp {
	case "subzerodegree":
		result.Subzerodegree += percentageCalculator(result.Subzerodegree*float64(change), percentage)
	case "zerodegree":
		result.Zerodegree += percentageCalculator(result.Zerodegree*float64(change), percentage)
	case "tendegree":
		result.Tendegree += percentageCalculator(result.Tendegree*float64(change), percentage)
	case "twentiedegree":
		result.Twentiedegree += percentageCalculator(result.Twentiedegree*float64(change), percentage)
	case "thirtydegree":
		result.Thirtydegree += percentageCalculator(result.Thirtydegree*float64(change), percentage)
	}

	arr := []float64{
		result.Mon,
		result.Tue,
		result.Wed,
		result.Thu,
		result.Fri,
		result.Sat,
		result.Sun,

		result.Win,
		result.Spr,
		result.Sum,
		result.Aut,

		result.Thirtydegree,
		result.Twentiedegree,
		result.Tendegree,
		result.Zerodegree,
		result.Subzerodegree,
	};

	result.Overall = calculateAverage(arr)
	fmt.Println(result.Overall)

	return result
}

func UpdateSelected(id string, change int, c *gin.Context) *gorm.DB {
	result := GetDataByID(id, c)

	result.Selected += change

	result.Rating = update(result.Rating, change, c)

	res := initializers.DB.Save(result.Rating)
	err := res.Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	return res
}
