package tools

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/models"
)

func CheckIfRecipeExists(title string) (bool, error) {
	var result struct {
		Found bool
	}

	err := initializers.DB.Raw("SELECT EXISTS(SELECT * FROM ingredients_schemas WHERE ingredient = ?) AS found;", title).Scan(&result).Error

	return result.Found, err
}

func GetIngredients(c *gin.Context) []string {
	var ingredients []models.IngredientsSchema
	data := GetCurrentData(c)

	// Fetch the first 100 IngredientsSchema records ordered by the sum of Mon, Win, and Subzerodegree
	result := initializers.DB.Table("ingredients_schemas").
		Select("ingredients_schemas.*").
		Joins("JOIN rating_structs ON rating_structs.owner_id = ingredients_schemas.id").
		Order("rating_structs." + data.Day + " + rating_structs." + data.Season + " + rating_structs." + data.Temp + " DESC").
		Preload(clause.Associations).
		Limit(100).
		Find(&ingredients)

	if result.Error != nil {
		// Handle the error
	}

	ingredientArray := make([]string, len(ingredients))

	// Extract the Ingredient field into the array
	for i, ingredient := range ingredients {
		ingredientArray[i] = ingredient.Ingredient
	}

	return ingredientArray
}
func GetRecipes(returnedIngredients []string) []models.RecipeSchema {
	var recipes []models.RecipeSchema

	// Fetch recipes that contain at least one of the returned ingredients
	result := initializers.DB.Table("recipe_schemas").
		Select("DISTINCT recipe_schemas.*").
		Joins("JOIN ingredients_schemas ON ingredients_schemas.recipe_schema_id = recipe_schemas.id").
		Where("ingredients_schemas.ingredient IN (?)", returnedIngredients).
		Preload(clause.Associations).
		Preload("Ingredients.Rating").
		Find(&recipes)

	if result.Error != nil {
		// Handle the error
	}

	fmt.Println()

	return recipes
}
