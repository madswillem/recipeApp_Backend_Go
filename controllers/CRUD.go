package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"rezeptapp.ml/goApp/error_handler"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/models"
)

func GetAll(c *gin.Context) {
	var recipes []models.RecipeSchema

	result := initializers.DB.Preload(clause.Associations).Preload("Ingredients.Rating").Preload("Ingredients.NutritionalValue").Find(&recipes)

	if result.Error != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Database error", result.Error)
	}

	c.JSON(http.StatusOK, recipes)
}

func AddRecipe(c *gin.Context) {
	var body models.RecipeSchema

	err := c.ShouldBindJSON(&body)

	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Failed to read body", err)
	}

	body.Rating.DefaultRatingStruct(body.Title)
	for i := 0; i < len(body.Ingredients); i++ {
		body.Ingredients[i].Rating.DefaultRatingStruct(body.Ingredients[i].Ingredient)
	}

	err = body.AddNutritionalValue(c)
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Database error", err)
		return
	}

	tx := initializers.DB.Begin()

	if err := tx.Create(&body).Error; err != nil {
		tx.Rollback()
		error_handler.HandleError(c, http.StatusBadRequest, "Database error", err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Database error", err)
		return
	}

	c.JSON(http.StatusCreated, body)
}

func GetById(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "id is not a number", err)
		return
	}
	response := models.RecipeSchema{ID: uint(i)}
	err = response.GetRecipeByID(c)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			error_handler.HandleError(c, http.StatusNotFound, "Recipe not found", err)
		} else {
			error_handler.HandleError(c, http.StatusInternalServerError, "Database error", err)
		}
		return
	}
	c.JSON(http.StatusOK, response)
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
		error_handler.HandleError(c, http.StatusBadRequest, "Failed to read body", err)
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
		error_handler.HandleError(c, http.StatusInternalServerError, "Database error", err)
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func Select(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "id is not a number", err)
		return
	}
	response := models.RecipeSchema{ID: uint(i)}
	err = response.UpdateSelected(1, c)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			error_handler.HandleError(c, http.StatusNotFound, "Recipe not found", err)
		} else {
			error_handler.HandleError(c, http.StatusInternalServerError, "Database error", err)
		}
		return
	}

	c.Status(http.StatusOK)
}

func Deselect(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "id is not a number", err)
		return
	}
	response := models.RecipeSchema{ID: uint(i)}
	err = response.UpdateSelected(-1, c)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			error_handler.HandleError(c, http.StatusNotFound, "Recipe not found", err)
		} else {
			error_handler.HandleError(c, http.StatusInternalServerError, "Database error", err)
		}
		return
	}

	c.Status(http.StatusOK)
}

func Colormode(c *gin.Context) {
	switch c.Param("type") {
	case "get":
		cookie, err := c.Cookie("type")
		if err != nil {
			error_handler.HandleError(c, http.StatusBadRequest, "Cookie error", err)
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