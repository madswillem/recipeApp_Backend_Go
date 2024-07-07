package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Server) RenderHome(c *gin.Context) {
	var recipes []models.RecipeSchema
	result := s.DB.Preload(clause.Associations).Preload("Ingredients.Rating").Preload("Ingredients.NutritionalValue").Find(&recipes)
	if result.Error != nil {
		panic(result.Error)
	}

	c.HTML(http.StatusOK, "construction/index", gin.H{
		"pageTitle": "Appetaizr",
		//		"recipes":   recipes,
	})
}
func (s *Server) RenderAcount(c *gin.Context) {
	c.HTML(http.StatusOK, "construction/index", gin.H{
		"pageTitle": "Appetaizr",
	})
}
func (s *Server) RenderTutorial(c *gin.Context) {
	c.HTML(http.StatusOK, "construction/index", gin.H{
		"pageTitle": "Appetaizr",
	})
}
func (s *Server) RenderProductpage(c *gin.Context) {
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
	getErr := res.GetRecipeByIDGORM(s.DB ,reqData)

	if getErr != nil {
		if getErr.Errors[0] == gorm.ErrRecordNotFound {
			error_handler.HandleError(c, getErr.Code, getErr.Message, getErr.Errors)
		} else {
			error_handler.HandleError(c, getErr.Code, getErr.Message, getErr.Errors)
		}
		return
	}

	c.HTML(http.StatusOK, "construction/index", gin.H{
		"pageTitle": "Appetaizr",
		//		"recipe":    res,
	})
}
