package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
	"gorm.io/gorm/clause"
)

func (s *Server) GetHome(c *gin.Context) {
	var recipes []models.RecipeSchema
	result := s.DB.Preload(clause.Associations).Preload("Ingredients.Rating").Find(&recipes)
	if result.Error != nil {
		panic(result.Error)
	}

	c.HTML(http.StatusOK, "construction/content", gin.H{
		"recipes": recipes,
	})
}
func (s *Server) GetRecipe(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "id is not a number", []error{err})
		return
	}
	res := models.RecipeSchema{}
	res.ID = fmt.Sprint(i)
	reqData := map[string]bool{
		"ingredients":      true,
		"ingredient_nutri": true,
		"ingredient_rate":  true,
		"rating":           true,
		"nutritionalvalue": true,
		"diet":             true,
	}
	res.GetRecipeByIDGORM(s.DB ,reqData)

	if res.ID == "" {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"pageTitle": "404 Page not found",
		})
		return
	}

	c.HTML(http.StatusOK, "construction/content", gin.H{
		"recipe": res,
	})
}
func (s *Server) GetAccount(c *gin.Context) {
	c.HTML(http.StatusOK, "construction/content", gin.H{})
}

func (s *Server) GetImgs(c *gin.Context) {
	name := c.Param("filename")
	fullName := filepath.Join(filepath.FromSlash(path.Clean("./web/imgs/" + name)))
	if _, err := os.Stat(fullName); errors.Is(err, os.ErrNotExist) {
		c.AbortWithStatus(http.StatusNotFound)
	}
	c.File(fullName)
}
func (s *Server) GetStyles(c *gin.Context) {
	name := c.Param("filename")
	fullName := filepath.Join(filepath.FromSlash(path.Clean("./web/styles/" + name)))
	if _, err := os.Stat(fullName); errors.Is(err, os.ErrNotExist) {
		c.AbortWithStatus(http.StatusNotFound)
	}
	c.File(fullName)
}
func (s *Server) GetScripts(c *gin.Context) {
	name := c.Param("filename")
	fullName := filepath.Join(filepath.FromSlash(path.Clean("./web/scripts/" + name)))
	if _, err := os.Stat(fullName); errors.Is(err, os.ErrNotExist) {
		c.AbortWithStatus(http.StatusNotFound)
	}
	c.File(fullName)
}
