package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/models"
)

func GetDataByID(id string, c *gin.Context)(models.RecipeSchema) {
	// find
	var result models.RecipeSchema
	err := initializers.DB.Preload("Ingredients").Find(&result, "ID = ?", id).Error

	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}

	fmt.Println(result)

	return result
}