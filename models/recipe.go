package models

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"rezeptapp.ml/goApp/error_handler"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/tools"
)

type RecipeSchema struct {
	ID        		 uint 			 	 `json:"_id" gorm:"primarykey"`
	Title			 string				 `json:"title"` 
	Ingredients		 []IngredientsSchema `json:"ingredients" gorm:"foreignKey:RecipeSchemaID"`
	Preparation		 string				 `json:"preparation"`
	CookingTime		 int				 `json:"cookingtime"`
	Image			 string				 `json:"image"`
	NutriScore		 string				 `json:"nutriscore"`
	NutritionalValue NutritionalValue    `json:"nutritional_value" gorm:"polymorphic:Owner"`
	Diet			 DietSchema			 `json:"diet" gorm:"foreignKey:RecipeID"`
	Selected		 int				 `json:"selected"`
	Rating			 RatingStruct		 `json:"rating" gorm:"polymorphic:Owner"`	
	Version     	 int				 `json:"__v"`
}

func (recipe *RecipeSchema) CheckIfExistsByTitle() (bool, error) {
	var result struct {
		Found bool
	}
	err := initializers.DB.Raw("SELECT EXISTS(SELECT * FROM ingredients_schemas WHERE ingredient = ?) AS found;", recipe.Title).Scan(&result).Error
	return result.Found, err
}

func (recipe *RecipeSchema) GetRecipeByID(c *gin.Context) error {
	err := initializers.DB.Preload(clause.Associations).Preload("Ingredients.NutritionalValue").Preload("Ingredients.Rating").Preload("Ingredients.NutritionalValue").First(&recipe, "ID = ?", recipe.ID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			error_handler.HandleError(c, http.StatusNotFound, "Recipe not found", err)
		} else {
			error_handler.HandleError(c, http.StatusInternalServerError, "Database error", err)
		}
	}

	return err
}

func (recipe *RecipeSchema) AddNutritionalValue(c *gin.Context) error {
	for _, ingredient := range recipe.Ingredients {
        var nutritionalValue NutritionalValue		
        err := initializers.DB.Joins("JOIN ingredients_schemas ON nutritional_values.owner_id = ingredients_schemas.id").
            Where("ingredients_schemas.ingredient = ?", ingredient.Ingredient).
            First(&nutritionalValue).Error
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) && !ingredient.NutritionalValue.Edited {
				print(ingredient.Ingredient)
				return err
            } else if errors.Is(err, gorm.ErrRecordNotFound) && ingredient.NutritionalValue.Edited {
				err = ingredient.createIngredientDBEntry(c)
				if err != nil {
					return err
				}
            } else {
				return err
            }
        } else if err == nil {
            if ingredient.NutritionalValue.Edited {
				err = errors.New("ingredient already exists")
				return err
            } else if !ingredient.NutritionalValue.Edited {
                ingredient.NutritionalValue = nutritionalValue
            }
        }
    }
	return nil
}

func (recipe *RecipeSchema) UpdateSelected(change int, c *gin.Context) error {
	recipe.Selected += change
	recipe.Rating.Update(change, c)

	res := initializers.DB.Save(recipe)
	err := res.Error
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Database error", err)
		return err
	}

	return err
}




func (rating RatingStruct) Update(change int, c *gin.Context) (RatingStruct, error) {

	result := rating
	data, err := tools.GetCurrentData()

	if err != nil {
		return result, err
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
		c.AbortWithStatus(http.StatusInternalServerError)
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
		c.AbortWithStatus(http.StatusInternalServerError)
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
	};

	result.Overall = tools.CalculateAverage(arr)
	fmt.Println(result.Overall)

	return result, err
}



func (ingredient *IngredientsSchema) createIngredientDBEntry(c *gin.Context) error {
	newIngredientDBEntry := IngredientDBSchema{
		Name:       ingredient.Ingredient,
		StandardUnit: ingredient.MeasurementUnit,
		NutritionalValue: ingredient.NutritionalValue,
	}
	
	err := initializers.DB.Create(&newIngredientDBEntry).Error
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Database error", err)
	}

	return err
}