package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/models"
)

func GetDataByID(id string, c *gin.Context)(models.RecipeSchema) {
	// find
	var result models.RecipeSchema
	err := initializers.DB.Find(&result, 10)

	if err.Error != nil {
		c.AbortWithError(http.StatusNotFound, err.Error)
	}

	return result
}