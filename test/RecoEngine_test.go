package test

import (
	"testing"

	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
)

func TestGetSimialrity(t *testing.T) {
	testGroup := models.RecipeGroupSchema{
		ID:      1,
		Recipes: []*models.RecipeSchema{{ID: 0}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 8}, {ID: 9}, {ID: 10}, {ID: 11}, {ID: 12}, {ID: 13}, {ID: 14}},
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
		ID: 1,
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
