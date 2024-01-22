package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/initializers"
	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetAll(c *gin.Context) {
	var recipes []models.RecipeSchema

	result := initializers.DB.Preload(clause.Associations).Preload("Ingredients.Rating").Preload("Ingredients.NutritionalValue").Find(&recipes)

	if result.Error != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Database error", []error{result.Error})
	}

	c.JSON(http.StatusOK, recipes)
}

func AddRecipe(c *gin.Context) {
	var body models.RecipeSchema

	binderr := c.ShouldBindJSON(&body)
	if binderr != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Failed to read body", []error{binderr})
		return
	}

	err := body.Create()
	if err != nil {
		error_handler.HandleError(c, err.Code, err.Message, err.Errors)
		return
	}

	c.JSON(http.StatusCreated, body)
}

func UpdateRecipe(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "id is not a number", []error{err})
		return
	}

	var body models.RecipeSchema
	c.ShouldBindJSON(&body)

	body.ID = uint(i)

	updateErr := body.Update()
	if updateErr != nil {
		error_handler.HandleError(c, updateErr.Code, updateErr.Message, updateErr.Errors)
		return
	}
}

func DeleteRecipe(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "id is not a number", []error{err})
		return
	}
	response := models.RecipeSchema{ID: uint(i)}
	deleteErr := response.Delete(c)
	if deleteErr != nil {
		error_handler.HandleError(c, deleteErr.Code, deleteErr.Message, deleteErr.Errors)
		return
	}
	c.Status(http.StatusOK)
}

func GetById(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errMessage := strings.Join([]string{"id", c.Param("id"), "is not a number"}, " ")
		error_handler.HandleError(c, http.StatusBadRequest, errMessage, []error{err})
		return
	}
	response := models.RecipeSchema{ID: uint(i)}
	reqData := map[string]bool{
		"ingredients":      true,
		"ingredient_nutri": true,
		"ingredient_rate":  true,
		"rating":           true,
		"nutritionalvalue": true,
		"diet":             true,
	}
	getErr := response.GetRecipeByID(reqData)

	if getErr != nil {
		if getErr.Errors[0] == gorm.ErrRecordNotFound {
			error_handler.HandleError(c, getErr.Code, getErr.Message, getErr.Errors)
		} else {
			error_handler.HandleError(c, getErr.Code, getErr.Message, getErr.Errors)
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
		error_handler.HandleError(c, http.StatusBadRequest, "Failed to read body", []error{err})
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
		error_handler.HandleError(c, http.StatusInternalServerError, "Database error", []error{err})
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func Select(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "id is not a number", []error{err})
		return
	}
	response := models.RecipeSchema{ID: uint(i)}
	selectedErr := response.UpdateSelected(1, c)
	if selectedErr != nil {
		error_handler.HandleError(c, selectedErr.Code, selectedErr.Message, selectedErr.Errors)
		return
	}

	c.Status(http.StatusOK)
}

func Deselect(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "id is not a number", []error{err})
		return
	}
	response := models.RecipeSchema{ID: uint(i)}
	selectedErr := response.UpdateSelected(-1, c)
	if selectedErr != nil {
		error_handler.HandleError(c, selectedErr.Code, selectedErr.Message, selectedErr.Errors)
		return
	}

	c.Status(http.StatusOK)
}

func Colormode(c *gin.Context) {
	switch c.Param("type") {
	case "get":
		cookie, err := c.Cookie("type")
		if err != nil {
			error_handler.HandleError(c, http.StatusBadRequest, "Cookie error", []error{err})
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
