package controllers

import (	
	"math/rand"
	"net/http"

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
	c.JSON(http.StatusOK, middleware.GetDataByID(c.Param("id"), c))
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

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Failed to read body",
			"errMessage": err.Error(),
		})
		return
	}

	var ingredientNames []string
	var recipes []models.RecipeSchema

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
		if err != nil {c.JSON(http.StatusBadRequest, gin.H{"err": err,})}
		c.JSON(http.StatusOK, gin.H{"type": cookie,})
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
	recipes := tools.GetRecipes(tools.GetIngredients(c))
	
	if len(recipes) <= 5 {
		c.AbortWithStatusJSON(http.StatusAccepted, recipes)
		return
	}

	var res [5]models.RecipeSchema
	for i := 0; i < 5; i++ {
		res[i] = recipes[rand.Intn(len(recipes) - 1)]
	}

	c.JSON(http.StatusAccepted, res)
}