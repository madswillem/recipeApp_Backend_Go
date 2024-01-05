package controllers

import (
	"errors"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/middleware"
	"rezeptapp.ml/goApp/models"
	"rezeptapp.ml/goApp/tools"
)

func handleError(c *gin.Context, statusCode int, errorMessage string, err error) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error":      errorMessage,
		"errMessage": err.Error(),
	})
	panic(errorMessage)
}

func GetAll(c *gin.Context) {
	var recipes []models.RecipeSchema

	result := initializers.DB.Preload(clause.Associations).Preload("Ingredients.Rating").Preload("Ingredients.NutritionalValue").Find(&recipes)

	if result.Error != nil {
		panic(result.Error)
	}

	c.JSON(http.StatusOK, recipes)
}

func AddRecipe(c *gin.Context) {
	var body models.RecipeSchema

	err := c.ShouldBindJSON(&body)

	if err != nil {
		handleError(c, http.StatusBadRequest, "Failed to read body", err)
	}

	body.Rating = *models.NewRatingStruct(body.Title)
	for i := 0; i < len(body.Ingredients); i++ {
		body.Ingredients[i].Rating = *models.NewRatingStruct(body.Ingredients[i].Ingredient)
	}

	body = CheckIfIngredientExists(c, body)

	tx := initializers.DB.Begin()

	if err := tx.Create(&body).Error; err != nil {
		tx.Rollback()
		handleError(c, http.StatusBadRequest, "Database error", err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		handleError(c, http.StatusBadRequest, "Database error", err)
		return
	}

	c.JSON(http.StatusCreated, body)
}

func CheckIfIngredientExists(c *gin.Context, body models.RecipeSchema) models.RecipeSchema{
	for _, ingredient := range body.Ingredients {
        var nutritionalValue models.NutritionalValue		
        err := initializers.DB.Joins("JOIN ingredients_schemas ON nutritional_values.owner_id = ingredients_schemas.id").
            Where("ingredients_schemas.ingredient = ?", ingredient.Ingredient).
            First(&nutritionalValue).Error

        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) && !ingredient.NutritionalValue.Edited {
                handleError(c, http.StatusBadRequest, "Database error", err)
				return models.RecipeSchema{}
            } else if errors.Is(err, gorm.ErrRecordNotFound) && ingredient.NutritionalValue.Edited {
				CreateIngredientDBEntry(c, ingredient)
            } else {
                handleError(c, http.StatusInternalServerError, "Database error", err)
				return models.RecipeSchema{}
            }
        } else if err == nil {
            if ingredient.NutritionalValue.Edited {
                handleError(c, http.StatusBadRequest, "Ingredient already exists", err)
				return models.RecipeSchema{}
            } else if !ingredient.NutritionalValue.Edited {
                ingredient.NutritionalValue = nutritionalValue
            }
        }
    }
	return body
}
func CreateIngredientDBEntry(c *gin.Context, ingredient models.IngredientsSchema) {
	newIngredientDBEntry := models.IngredientDBSchema{
		Name:       ingredient.Ingredient,
		StandardUnit: ingredient.MeasurementUnit,
		NutritionalValue: ingredient.NutritionalValue,
	}
	
	if err := initializers.DB.Create(&newIngredientDBEntry).Error; err != nil {
		handleError(c, http.StatusBadRequest, "Database error", err)
		return	
	}
}

func GetById(c *gin.Context) {
	c.JSON(http.StatusOK, middleware.GetDataByID(c.Param("id"), c))
}

func Filter(c *gin.Context) {
	type Recipe struct {
		NutriScore  string            `json:"nutriscore"`
		CookingTime int               `json:"cookingtime"`
		Ingredients []string          `json:"ingredients"`
		Diet        models.DietSchema `json:"diet"`
	}

	var body Recipe
	err := c.ShouldBindJSON(&body)

	if err != nil {
		handleError(c, http.StatusBadRequest, "Failed to read body", err)
		return
	}

	var recipes []models.RecipeSchema

	query := initializers.DB.Joins("JOIN ingredients_schemas ON recipe_schemas.id = ingredients_schemas.recipe_schema_id").
		Where("ingredients_schemas.ingredient IN ?", body.Ingredients).
		Group("recipe_schemas.id").
		Having("COUNT(DISTINCT ingredients_schemas.id) = ?", len(body.Ingredients)).
		Joins("JOIN diet_schemas ON diet_schemas.recipe_id = recipe_schemas.id").
		Preload(clause.Associations).
		Preload("Ingredients.Rating").
		Preload("Ingredients.NutritionalValue")

		switch {
		case body.Diet.Vegetarien:
				query = query.Where("diet_schemas.vegetarien = ?", true)
		case body.Diet.Vegan:
			query = query.Where("diet_schemas.vegan = ?", true)
		case body.Diet.LowCal:
			query = query.Where("diet_schemas.lowcal = ?", true)
		case body.Diet.LowCarb:
			query = query.Where("diet_schemas.lowcarb = ?", true)
		case body.Diet.Keto:
			query = query.Where("diet_schemas.keto = ?", true)
		case body.Diet.Paleo:
			query = query.Where("diet_schemas.paleo = ?", true)
		case body.Diet.LowFat:
			query = query.Where("diet_schemas.lowfat = ?", true)
		case body.Diet.FoodCombining:
			query = query.Where("diet_schemas.food_combining = ?", true)
		case body.Diet.WholeFood:
			query = query.Where("diet_schemas.whole_food = ?", true)
		case body.CookingTime > 0:
			query = query.Where("recipe_schemas.cooking_time <= ?", body.CookingTime)
		case body.NutriScore != "":
			query = query.Where("recipe_schemas.nutri_score = ?", body.NutriScore)
		}

	err = query.Find(&recipes).Error

	if err != nil {
		handleError(c, http.StatusInternalServerError, "Database error", err)
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func Select(c *gin.Context) {
	c.JSON(http.StatusOK, middleware.UpdateSelected(c.Param("id"), +1, c).Error)
}

func Deselect(c *gin.Context) {
	res := middleware.UpdateSelected(c.Param("id"), -1, c)

	c.JSON(http.StatusOK, res.Error)
}

func Colormode(c *gin.Context) {
	switch c.Param("type") {
	case "get":
		cookie, err := c.Cookie("type")
		if err != nil {
			handleError(c, http.StatusBadRequest, "Cookie error", err)
		}
		c.JSON(http.StatusOK, gin.H{"type": cookie})
	case "dark":
		c.SetCookie("type", "dark", 999999999999999999, "/", "localhost", false, true)
		c.Status(http.StatusAccepted)
	case "light":
		c.SetCookie("type", "light", 999999999999999999, "/", "localhost", false, true)
		c.Status(http.StatusAccepted)
	default:
		c.Status(http.StatusBadRequest)
	}
}

func Recomend(c *gin.Context) {
	recipes := tools.GetRecipes(c, tools.GetIngredients(c))

	if len(recipes) <= 5 {
		c.AbortWithStatusJSON(http.StatusAccepted, recipes)
		return
	}

	var res [5]models.RecipeSchema
	for i := 0; i < 5; i++ {
		res[i] = recipes[rand.Intn(len(recipes)-1)]
	}

	c.JSON(http.StatusAccepted, res)
}