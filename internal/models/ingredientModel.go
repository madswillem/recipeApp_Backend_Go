package models

import (
	"errors"
	"net/http"

	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/initializers"
)

type IngredientsSchema struct {
	ID        		 uint 			 	 `json:"id" gorm:"primarykey"`
	RecipeSchemaID 	 uint				 `json:"-"`

	Ingredient		 string				 `json:"ingredient"` 
	Amount			 string				 `json:"amount"` 
	MeasurementUnit	 string				 `json:"measurement_unit"`
	NutritionalValue NutritionalValue    `json:"nutritional_value" gorm:"polymorphic:Owner"`
	Rating			 RatingStruct		 `json:"rating" gorm:"polymorphic:Owner"`
}

func (ingredient *IngredientsSchema) createIngredientDBEntry() *error_handler.APIError {
	newIngredientDBEntry := IngredientDBSchema{
		Name:       ingredient.Ingredient,
		StandardUnit: ingredient.MeasurementUnit,
		NutritionalValue: ingredient.NutritionalValue,
	}
	
	err := initializers.DB.Create(&newIngredientDBEntry).Error
	if err != nil {
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}

	return nil
}

func (ingredient *IngredientsSchema) CheckForRequiredFields() error {
	if ingredient.Ingredient == "" {
		return errors.New("missing ingredient")
	}
	if ingredient.Amount == "" {
		return errors.New("missing amount")
	}
	if ingredient.MeasurementUnit == "" {
		return errors.New("missing measurement unit")
	}
	return nil
}
