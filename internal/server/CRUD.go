package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Server) GetAll(c *gin.Context) {
	var recipes []models.RecipeSchema
	result := s.DB.Preload(clause.Associations).Preload("Ingredients.Rating").Preload("Ingredients.NutritionalValue").Find(&recipes)

	if result.Error != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Database error", []error{result.Error})
	}
	c.JSON(http.StatusOK, recipes)
}

func (s *Server) AddRecipe(c *gin.Context) {
	var body models.RecipeSchema

	binderr := c.ShouldBindJSON(&body)
	if binderr != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Failed to read body", []error{binderr})
		return
	}

	err := body.Create(s.DB)
	if err != nil {
		error_handler.HandleError(c, err.Code, err.Message, err.Errors)
		return
	}

	c.JSON(http.StatusCreated, body)
}

func (s *Server) UpdateRecipe(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "id is not a number", []error{err})
		return
	}

	var body models.RecipeSchema
	c.ShouldBindJSON(&body)

	body.ID = uint(i)

	updateErr := body.Update(s.DB)
	if updateErr != nil {
		error_handler.HandleError(c, updateErr.Code, updateErr.Message, updateErr.Errors)
		return
	}
}

func (s *Server) DeleteRecipe(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "id is not a number", []error{err})
		return
	}
	response := models.RecipeSchema{}
	response.ID = uint(i)

	deleteErr := response.Delete(s.DB)
	if deleteErr != nil {
		error_handler.HandleError(c, deleteErr.Code, deleteErr.Message, deleteErr.Errors)
		return
	}
	c.Status(http.StatusOK)
}

func (s *Server) GetById(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errMessage := strings.Join([]string{"id", c.Param("id"), "is not a number"}, " ")
		error_handler.HandleError(c, http.StatusBadRequest, errMessage, []error{err})
		return
	}
	response := models.RecipeSchema{}
	response.ID = uint(i)

	reqData := map[string]bool{
		"ingredients":      true,
		"ingredient_nutri": true,
		"ingredient_rate":  true,
		"rating":           true,
		"nutritionalvalue": true,
		"diet":             true,
	}
	getErr := response.GetRecipeByID(s.DB ,reqData)

	if getErr != nil {
		if getErr.Errors[0] == gorm.ErrRecordNotFound {
			error_handler.HandleError(c, getErr.Code, getErr.Message, getErr.Errors)
		} else {
			error_handler.HandleError(c, getErr.Code, getErr.Message, getErr.Errors)
		}
		return
	}
	c.JSON(http.StatusOK, response)
}

func (s *Server) Filter(c *gin.Context) {

	var body models.Filter
	err := c.ShouldBindJSON(&body)

	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Failed to read body", []error{err})
		return
	}

	recipes, apiErr := body.Filter()
	if apiErr != nil {
		error_handler.HandleError(c, apiErr.Code, apiErr.Message, apiErr.Errors)
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func (s *Server) Select(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "id is not a number", []error{err})
		return
	}
	middleware_user, _ := c.MustGet("user").(models.UserModel)
	
	response := models.RecipeSchema{}
	response.ID = uint(i)
	
	selectedErr := response.UpdateSelected(1, &middleware_user, s.DB)
	if selectedErr != nil {
		error_handler.HandleError(c, selectedErr.Code, selectedErr.Message, selectedErr.Errors)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) Deselect(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "id is not a number", []error{err})
		return
	}
	middleware_user, _ := c.MustGet("user").(models.UserModel)

	response := models.RecipeSchema{}
	response.ID = uint(i)

	selectedErr := response.UpdateSelected(-1, &middleware_user, s.DB)
	if selectedErr != nil {
		error_handler.HandleError(c, selectedErr.Code, selectedErr.Message, selectedErr.Errors)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) Colormode(c *gin.Context) {
	switch c.Param("type") {
	case "get":
		cookie, err := c.Cookie("type")
		if err != nil {
			error_handler.HandleError(c, http.StatusBadRequest, "Cookie error", []error{err})
		}
		c.JSON(http.StatusOK, gin.H{"type": cookie})
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
