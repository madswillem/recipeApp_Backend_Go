package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/middleware"
	"rezeptapp.ml/goApp/models"
	"rezeptapp.ml/goApp/tools"
)

func GetAll(c *gin.Context) {
	var recipes []models.RecipeSchema

	result := initializers.DB.Preload("Ingredients").Find(&recipes)

	if result.Error != nil {
		panic(result.Error)
	}

	c.JSON(http.StatusOK, recipes)
}

func AddRecipe(c *gin.Context) {
	var body models.RecipeSchema
	id := tools.NewObjectId()

	err := c.Bind(&body)
	var recipes = make([]models.IngredientsSchema, len(body.Ingredients))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Failed to read body",
			"errMessage": err.Error(),
		})
		return
	}

	body.ID = id
	for i:=0; i<len(body.Ingredients); i++ { 
		result, err := tools.CheckIfRecipeExists(body.Ingredients[i].Ingredient)  
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":      "Failed to read body",
				"errMessage": err.Error(),
			})
			return
		}
		if result {
			err = initializers.DB.Find(&recipes[i], "ingredient = ?", body.Ingredients[i].Ingredient).Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":      "Failed to read body",
					"errMessage": err.Error(),
				})
				return
			}
		}
		if !result {
			ingId := tools.NewObjectId()          // start of the execution block
        	body.Ingredients[i].ID = ingId
			recipes[i] = body.Ingredients[i]
		}	
    } 

	body.Ingredients = recipes

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
		Ingredient		string			`json:"ingredient"`
	} 
	type Recipe struct {
		Ingredients		[]Ingredients	`json:"ingredients"`
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

	err = initializers.DB.Joins("JOIN ingredients_schemas ON recipe_schemas.id = ingredients_schemas.recipe_schema_id").
	Where("ingredients_schemas.ingredient IN ?", ingredientNames).
	Group("recipe_schemas.id").
	Having("COUNT(DISTINCT ingredients_schemas.id) = ?", len(ingredientNames)).
	Preload("Ingredients").
	Find(&recipes).Error

	if err != nil {
		panic(err.Error)
	}
	if len(recipes) <= 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, "No Recipe was Found")
		return
	}

	c.JSON(http.StatusOK , recipes)
}

func Select(c *gin.Context) {
	res := middleware.UpdateSelected(c.Param("id"), +1, c)

	c.JSON(http.StatusOK, res)
}

func Deselect(c *gin.Context) {
	res := middleware.UpdateSelected(c.Param("id"), -1, c)

	c.JSON(http.StatusOK, res)
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