package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/models"
)

func GetHome(c *gin.Context) {
	var recipes []models.RecipeSchema
	result := initializers.DB.Preload(clause.Associations).Preload("Ingredients.Rating").Find(&recipes)
	if result.Error != nil {panic(result.Error)}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"recipes": recipes,
	})
}
func GetAccount(c *gin.Context) {
	c.HTML(http.StatusOK, "account.html", gin.H{})
}