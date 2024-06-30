package models

import (
	"errors"
	"net/http"
	"time"

	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/tools"
)

type RatingStruct struct {
	ID             string    `db:"id" json:"id"`
    CreatedAt      time.Time `db:"created_at" json:"created_at"`
    IngredientID   *string   `db:"ingredient_id" json:"ingredient_id,omitempty"`
    RecipeID       *string   `db:"recipe_id" json:"recipe_id,omitempty"`
    Overall        float64   `db:"overall" json:"overall"`
    Mon            float64   `db:"mon" json:"mon"`
    Tue            float64   `db:"tue" json:"tue"`
    Wed            float64   `db:"wed" json:"wed"`
    Thu            float64   `db:"thu" json:"thu"`
    Fri            float64   `db:"fri" json:"fri"`
    Sat            float64   `db:"sat" json:"sat"`
    Sun            float64   `db:"sun" json:"sun"`
    Win            float64   `db:"win" json:"win"`
    Spr            float64   `db:"spr" json:"spr"`
    Sum            float64   `db:"sum" json:"sum"`
    Aut            float64   `db:"aut" json:"aut"`
    ThirtyDegree   float64   `db:"thirtydegree" json:"thirtydegree"`
    TwentyDegree   float64   `db:"twentiedegree" json:"twentiedegree"`
    TenDegree      float64   `db:"tendegree" json:"tendegree"`
    ZeroDegree     float64   `db:"zerodegree" json:"zerodegree"`
    SubZeroDegree  float64   `db:"subzerodegree" json:"subzerodegree"`
}

func (rating *RatingStruct) DefaultRatingStruct(recipe_id *string, ingredient_id *string) {
	rating.RecipeID = recipe_id
	rating.IngredientID = ingredient_id

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
	rating.ThirtyDegree = 1000.0
	rating.TwentyDegree = 1000.0
	rating.TenDegree = 1000.0
	rating.ZeroDegree = 1000.0
	rating.SubZeroDegree = 1000.0
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
		result.SubZeroDegree += tools.PercentageCalculator(result.SubZeroDegree*float64(change), percentage)
	case "zerodegree":
		result.ZeroDegree += tools.PercentageCalculator(result.ZeroDegree*float64(change), percentage)
	case "tendegree":
		result.TenDegree += tools.PercentageCalculator(result.TenDegree*float64(change), percentage)
	case "twentiedegree":
		result.TwentyDegree += tools.PercentageCalculator(result.TwentyDegree*float64(change), percentage)
	case "thirtydegree":
		result.ThirtyDegree += tools.PercentageCalculator(result.ThirtyDegree*float64(change), percentage)
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

		result.ThirtyDegree,
		result.TwentyDegree,
		result.TenDegree,
		result.ZeroDegree,
		result.SubZeroDegree,
	}

	result.Overall = tools.CalculateAverage(arr)

	return nil
}
