package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madswillem/recipeApp_Backend_Go/internal/database"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
)

func (s *Server) UpadateUser(c *gin.Context) {
	var settings models.UserSettings
	err := c.ShouldBind(&settings)
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Couldn't read body", []error{err})
		return
	}
	middleware_user, _ := c.MustGet("user").(models.UserModel)
	apiErr := middleware_user.GetByCookie(s.NewDB)
	if apiErr != nil {
		error_handler.HandleError(c, apiErr.Code, apiErr.Message, apiErr.Errors)
		return
	}
	middleware_user.Settings = settings
	apiErr = database.Update(s.DB, &middleware_user)
	if apiErr != nil {
		error_handler.HandleError(c, apiErr.Code, apiErr.Message, apiErr.Errors)
		return
	}

	c.JSON(http.StatusAccepted, middleware_user)
}
func (s *Server) GetRecommendation(c *gin.Context) {
	middleware_user, _ := c.Get("user")
	user, ok := middleware_user.(models.UserModel)
	if !ok {
		fmt.Println("type assertion failed")
	}

	err := user.GetByCookie(s.NewDB)
	if err != nil {
		error_handler.HandleError(c, err.Code, err.Message, err.Errors)
		return
	}

	err, recipes := user.GetRecomendation(s.DB)
	if err != nil {
		error_handler.HandleError(c, err.Code, err.Message, err.Errors)
		return
	}

	c.JSON(http.StatusOK, recipes)
}
func (s *Server) CreateGroup(c *gin.Context) {
	middleware_user, _ := c.Get("user")
	user, ok := middleware_user.(models.UserModel)
	if !ok {
		fmt.Println("type assertion failed")
	}

	r := models.RecipeSchema{ID: "aa85daf1-dbc5-462d-a6fe-3fbb358b08dd"}
	apiErr := r.GetRecipeByID(s.NewDB)
	if apiErr != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, apiErr)
		return
	}

	rp := models.RecipeGroupSchema{}
	rp.Create(&r)
	user.RecipeGroups = append(user.RecipeGroups, rp)

	v, err := json.Marshal(user.RecipeGroups)
	if err != nil {
		error_handler.HandleError(c, http.StatusInternalServerError, "Couldnt Marshal recipe group", []error{err})
	}

	s.NewDB.MustExec(`UPDATE "user" SET groups = $1 WHERE id = $2`, v, user.ID)
	c.JSON(http.StatusAccepted, user)
}
