package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/middleware"
	"rezeptapp.ml/goApp/models"
	"rezeptapp.ml/goApp/tools"
)

func GetAll(c *gin.Context) {
	var recipes []models.RecipeSchema

	result := initializers.DB.Preload(clause.Associations).Preload("Ingredients.Rating").Find(&recipes)

	if result.Error != nil {
		panic(result.Error)
	}

	c.JSON(http.StatusOK, recipes)
}

func AddRecipe(c *gin.Context) {
	var body models.RecipeSchema

	err := c.Bind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Failed to read body",
			"errMessage": err.Error(),
		})
		return
	}

	body.ID = tools.NewObjectId()
	body.Rating = *models.NewRatingStruct(tools.NewObjectId(), body.Title)

	for i := 0; i < len(body.Ingredients); i++ {
		body.Ingredients[i].ID = tools.NewObjectId()
		body.Ingredients[i].Rating = *models.NewRatingStruct(tools.NewObjectId(), body.Ingredients[i].Ingredient)
	}

	fmt.Println(body.Rating)

	initializers.DB.Create(&body)
	result := initializers.DB.Save(&body)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Databaseerror",
			"errMessage": result.Error,
		})
		return
	}

	c.JSON(http.StatusCreated, body)
}

func GetById(c *gin.Context) {
	result := middleware.GetDataByID(c.Param("id"), c)

	c.JSON(http.StatusOK, result)
}

func Filter(c *gin.Context) {
	type Ingredients struct {
		Ingredient string `json:"ingredient"`
	}
	type Recipe struct {
		CookingTime int           `json:"cookingtime"`
		Ingredients []Ingredients `json:"ingredients"`
	}

	var body Recipe
	err := c.Bind(&body)

	// find
	var ingredientNames []string
	var recipes []models.RecipeSchema

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Failed to read body",
			"errMessage": err.Error(),
		})
		return
	}

	for i := 0; i < len(body.Ingredients); i++ {
		ingredientNames = append(ingredientNames, body.Ingredients[i].Ingredient)
	}

	query := initializers.DB.Joins("JOIN ingredients_schemas ON recipe_schemas.id = ingredients_schemas.recipe_schema_id").
		Where("ingredients_schemas.ingredient IN ?", ingredientNames).
		Group("recipe_schemas.id").
		Having("COUNT(DISTINCT ingredients_schemas.id) = ?", len(ingredientNames)).
		Preload("Ingredients")

	if body.CookingTime > 0 {
		query = query.Where("recipe_schemas.cooking_time <= ?", body.CookingTime)
	}

	err = query.Find(&recipes).Error

	if err != nil {
		panic(err.Error)
	}
	if len(recipes) <= 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, "No Recipe was Found")
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func Select(c *gin.Context) {
	res := middleware.UpdateSelected(c.Param("id"), +1, c)

	c.JSON(http.StatusOK, res.Error)
}

func Deselect(c *gin.Context) {
	res := middleware.UpdateSelected(c.Param("id"), -1, c)

	c.JSON(http.StatusOK, res.Error)
}

func Colormode(c *gin.Context) {
	if c.Param("type") == "get" {
		cookie, err := c.Cookie("type")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"type": cookie,
		})
	} else if c.Param("type") == "dark" {
		c.SetCookie("type", "dark", 999999999999999999, "/", "localhost", false, true)
		c.Status(http.StatusAccepted)
	} else if c.Param("type") == "light" {
		c.SetCookie("type", "light", 999999999999999999, "/", "localhost", false, true)
		c.Status(http.StatusAccepted)
	} else {
		c.Status(http.StatusBadRequest)
	}
}

func Recomend(c *gin.Context) {
	recipes := tools.GetRecipes(tools.GetIngredients(c))
	res := make([]models.RecipeSchema, 5)

	rand.Seed(time.Now().UnixNano())

	if len(recipes) < 5 {
		c.JSON(http.StatusAccepted, recipes)
	}

	for i := 0; i < 5; i++ {
		randomNumber := rand.Intn(len(recipes) - 1)
		res[i] = recipes[randomNumber]
	}

	c.JSON(http.StatusAccepted, res)
}