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
	"rezeptapp.ml/goApp/middleware"
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
func GetRecipe(c *gin.Context) {
	res := middleware.GetDataByID(c.Param("id"), c)

	if res.ID == 0 {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"pageTitle": "404 Page not found",
		})
		return
	}

	c.HTML(http.StatusOK, "productpage.html", gin.H{
		"recipe": res,
	})
}
func GetAccount(c *gin.Context) {
	c.HTML(http.StatusOK, "account.html", gin.H{})
}
func GetTutorials(c *gin.Context)  {
	c.HTML(http.StatusOK, "tutorials.html", gin.H{})
}


func GetImgs(c *gin.Context)  {
	name := c.Param("filename")
	fullName := filepath.Join(filepath.FromSlash(path.Clean("./imgs/"+name)))
	if _, err := os.Stat(fullName); errors.Is(err, os.ErrNotExist) {
		c.AbortWithStatus(http.StatusNotFound)
	}
    c.File(fullName)
}
func GetVideos(c *gin.Context) {
	videoName := c.Param("filename")
	c.Header("Content-Type", "application/vnd.apple.mpegurl")
	fullName := filepath.Join(filepath.FromSlash(path.Clean("./videos/out/" + videoName )))
	c.File(fullName)
}
func GetStyles(c *gin.Context)  {
	name := c.Param("filename")
	fullName := filepath.Join(filepath.FromSlash(path.Clean("./styles/"+name)))
	if _, err := os.Stat(fullName); errors.Is(err, os.ErrNotExist) {
		c.AbortWithStatus(http.StatusNotFound)
	}
    c.File(fullName)
}