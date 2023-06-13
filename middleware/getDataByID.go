package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/models"
)

func GetDataByID(id string, c *gin.Context) models.RecipeSchema {
	// find
	var result models.RecipeSchema
	err := initializers.DB.Preload(clause.Associations).Preload("Ingredients.Rating").Find(&result, "ID = ?", id).Error

	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}

	return result
}
