package controllers

import (
	"errors"

	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/models"
)

func RenderHome(c *gin.Context) {
	var recipes []models.RecipeSchema
	result := initializers.DB.Preload(clause.Associations).Preload("Ingredients.Rating").Find(&recipes)
	if result.Error != nil {panic(result.Error)}

	c.HTML(http.StatusOK, "home/index", gin.H{
		"pageTitle": "Recipe App",
		"recipes": recipes,
	})
}
func RenderAcount(c *gin.Context) {
	c.HTML(http.StatusOK, "account/index", gin.H{
		"title": "Recipe App",
	})
}

func GetRescources(c *gin.Context)  {
	name := c.Param("filename")
	fullName := filepath.Join(filepath.FromSlash(path.Clean("./rescources/"+name)))
	if _, err := os.Stat(fullName); errors.Is(err, os.ErrNotExist) {
		c.AbortWithStatus(http.StatusNotFound)
	}
    c.File(fullName)
}

func GetImgs(c *gin.Context)  {
	name := c.Param("filename")
	fullName := filepath.Join(filepath.FromSlash(path.Clean("./imgs/"+name)))
	if _, err := os.Stat(fullName); errors.Is(err, os.ErrNotExist) {
		c.AbortWithStatus(http.StatusNotFound)
	}
    c.File(fullName)
}
