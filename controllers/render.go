package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/middleware"
	"rezeptapp.ml/goApp/models"
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
	res := middleware.GetDataByID(c.Param("id"), c)

	if res.ID == "" {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"pageTitle": "404 Page not found",
		})
		return
	}

	c.HTML(http.StatusOK, "productpage/index", gin.H{
		"pageTitle": "Appetaizr",
		"recipe": res,
	})
}
