package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
)

func UpadateUser(c *gin.Context) {
	var settings models.UserSettings
	err := c.ShouldBind(&settings)
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Couldn't read body", []error{err})
		return
	}
	middleware_user, _ := c.MustGet("user").(models.UserModel)
	apiErr := middleware_user.GetByCookie()
	if apiErr != nil {
		error_handler.HandleError(c, apiErr.Code, apiErr.Message, apiErr.Errors)
		return
	}
	middleware_user.Settings = settings
	apiErr = middleware_user.Update()
	if apiErr != nil {
		error_handler.HandleError(c, apiErr.Code, apiErr.Message, apiErr.Errors)
		return
	}
	
	c.JSON(http.StatusAccepted, middleware_user)
}
