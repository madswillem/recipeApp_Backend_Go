package test

import (
	"testing"

	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
	"gorm.io/gorm"
)

func TestGetSimialrity(t *testing.T) {
	testGroup := models.RecipeGroupSchema{
		Model: gorm.Model{ID: 1},
		Recipes:[]*models.RecipeSchema{
			{Model: gorm.Model{ID: 1}},
			{Model: gorm.Model{ID: 2}},
			{Model: gorm.Model{ID: 3}},
			{Model: gorm.Model{ID: 4}},
			{Model: gorm.Model{ID: 5}},
			{Model: gorm.Model{ID: 6}},
			{Model: gorm.Model{ID: 7}},
			{Model: gorm.Model{ID: 8}},
			{Model: gorm.Model{ID: 9}},
			{Model: gorm.Model{ID: 10}},
			{Model: gorm.Model{ID: 11}},
			{Model: gorm.Model{ID: 12}},
			{Model: gorm.Model{ID: 13}},
			{Model: gorm.Model{ID: 14}},
		},
		AvrgIngredients: []models.Avrg{
			{Name: "flower", Percentige: 14},
			{Name: "water", Percentige: 14},
			{Name: "salt", Percentige: 14},
			{Name: "oliveoil", Percentige: 14},
			{Name: "yeast", Percentige: 14},
			{Name: "sugar", Percentige: 14},
		},
		AvrgCuisine: []models.Avrg{
			{Name: "asian", Percentige: 14},
		},
		AvrgVegetarien:    0.0,
		AvrgVegan:         0.0,
		AvrgLowCal:        0.0,
		AvrgLowCarb:       0.0,
		AvrgKeto:          0.0,
		AvrgPaleo:         0.0,
		AvrgLowFat:        0.0,
		AvrgFoodCombining: 0.0,
		AvrgWholeFood:     0.0,
	}

	testRecipe := models.RecipeSchema{
		Model: gorm.Model{ID: 1},
		Ingredients: []models.IngredientsSchema{
			{Ingredient: "flower"},
			{Ingredient: "water"},
			{Ingredient: "salt"},
			{Ingredient: "oliveoil"},
			{Ingredient: "yeast"},
		},
		Cuisine: "asian",
	}

	res, err := testRecipe.GetSimilarityWithGroup(testGroup)
	if err != nil {
		t.Errorf("Method returned error: %s", err.Message)
	}
	if res != 1.0 {
		t.Errorf("Expected %f but got %f", 1.0, res)
	}
}
