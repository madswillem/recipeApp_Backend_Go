package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/initializers"
	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
)

func RenderHome(c *gin.Context) {
	var recipes []models.RecipeSchema
	result := initializers.DB.Preload(clause.Associations).Preload("Ingredients.Rating").Preload("Ingredients.NutritionalValue").Find(&recipes)
	if result.Error != nil {panic(result.Error)}

	c.HTML(http.StatusOK, "home/index", gin.H{
		"pageTitle": "Appetaizr",
		"recipes": recipes,
	})
}
func RenderAcount(c *gin.Context) {
	c.HTML(http.StatusOK, "account/index", gin.H{
		"pageTitle": "Appetaizr",
	})
}
func RenderTutorial(c *gin.Context) {
	c.HTML(http.StatusOK, "tutorials/index", gin.H{
		"pageTitle": "Appetaizr",
	})
}
func RenderProductpage(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "id is not a number", err)
		return
	}
	
	res := models.RecipeSchema{ID: uint(i)}
	reqData := map[string]bool{
		"ingredients":      true,
		"ingredient_nutri": true,
		"ingredient_rate":  true,
		"rating":           true,
		"nutritionalvalue": true,
		"diet":             true,
	}
	err = res.GetRecipeByID(c, reqData)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			error_handler.HandleError(c, http.StatusNotFound, "Recipe not found", err)
		} else {
			error_handler.HandleError(c, http.StatusInternalServerError, "Database error", err)
		}
		return
	}

	c.HTML(http.StatusOK, "productpage/index", gin.H{
		"pageTitle": "Appetaizr",
		"recipe": res,
	})
}
