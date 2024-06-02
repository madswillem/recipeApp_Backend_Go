package models

import (
	"errors"
	"net/http"

	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/tools"
)

type RatingStruct struct {
	ID         uint   `gorm:"primarykey"`
	OwnerTitle string `json:"owner_title"`
	OwnerID    string
	OwnerType  string

	Overall float64 `json:"overall"`

	Mon float64 `json:"mon"`
	Tue float64 `json:"tue"`
	Wed float64 `json:"wed"`
	Thu float64 `json:"thu"`
	Fri float64 `json:"fri"`
	Sat float64 `json:"sat"`
	Sun float64 `json:"sun"`

	Win float64 `json:"win"`
	Spr float64 `json:"spr"`
	Sum float64 `json:"sum"`
	Aut float64 `json:"aut"`

	Thirtydegree  float64 `json:"thirtydegree"`
	Twentiedegree float64 `json:"twentiedegree"`
	Tendegree     float64 `json:"tendegree"`
	Zerodegree    float64 `json:"zerodegree"`
	Subzerodegree float64 `json:"subzerodegree"`
}

func (rating *RatingStruct) DefaultRatingStruct(title string) {
	rating.OwnerTitle = title

	rating.Overall = 1000.0
	rating.Mon = 1000.0
	rating.Tue = 1000.0
	rating.Wed = 1000.0
	rating.Thu = 1000.0
	rating.Fri = 1000.0
	rating.Sat = 1000.0
	rating.Sun = 1000.0
	rating.Win = 1000.0
	rating.Spr = 1000.0
	rating.Sum = 1000.0
	rating.Aut = 1000.0
	rating.Thirtydegree = 1000.0
	rating.Twentiedegree = 1000.0
	rating.Tendegree = 1000.0
	rating.Zerodegree = 1000.0
	rating.Subzerodegree = 1000.0
}

func (rating *RatingStruct) Update(change int) *error_handler.APIError {

	result := rating
	data, err := tools.GetCurrentData()

	if err != nil {
		return error_handler.New("fatal internal error", http.StatusInternalServerError, err)
	}

	percentage := 10.0

	switch data.Day {
	case "Mon":
		result.Mon += tools.PercentageCalculator(result.Mon*float64(change), percentage)
	case "Tue":
		result.Tue += tools.PercentageCalculator(result.Tue*float64(change), percentage)
	case "Wed":
		result.Wed += tools.PercentageCalculator(result.Wed*float64(change), percentage)
	case "Thu":
		result.Thu += tools.PercentageCalculator(result.Thu*float64(change), percentage)
	case "Fri":
		result.Fri += tools.PercentageCalculator(result.Fri*float64(change), percentage)
	case "Sat":
		result.Sat += tools.PercentageCalculator(result.Sat*float64(change), percentage)
	case "Sun":
		result.Sun += tools.PercentageCalculator(result.Sun*float64(change), percentage)
	default:
		return error_handler.New("fatal internal error", http.StatusInternalServerError, errors.New("fatal internal error"))
	}

	switch data.Season {
	case "Win":
		result.Win += tools.PercentageCalculator(result.Win*float64(change), percentage)
	case "Spr":
		result.Spr += tools.PercentageCalculator(result.Spr*float64(change), percentage)
	case "Sum":
		result.Sum += tools.PercentageCalculator(result.Sum*float64(change), percentage)
	case "Aut":
		result.Aut += tools.PercentageCalculator(result.Aut*float64(change), percentage)
	default:
		return error_handler.New("fatal internal error", http.StatusInternalServerError, errors.New("fatal internal error"))
	}

	switch data.Temp {
	case "subzerodegree":
		result.Subzerodegree += tools.PercentageCalculator(result.Subzerodegree*float64(change), percentage)
	case "zerodegree":
		result.Zerodegree += tools.PercentageCalculator(result.Zerodegree*float64(change), percentage)
	case "tendegree":
		result.Tendegree += tools.PercentageCalculator(result.Tendegree*float64(change), percentage)
	case "twentiedegree":
		result.Twentiedegree += tools.PercentageCalculator(result.Twentiedegree*float64(change), percentage)
	case "thirtydegree":
		result.Thirtydegree += tools.PercentageCalculator(result.Thirtydegree*float64(change), percentage)
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
	}

	result.Overall = tools.CalculateAverage(arr)

	return nil
}
