package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/models"
)

func UpdateSelected(id string, change int, c *gin.Context)(*gorm.DB) {

	result := GetDataByID(id, c)

	result.Version += 1

	res := initializers.DB.Model(&result).Updates(models.RecipeSchema{Selected: result.Selected + change, Version: result.Version +1})
	if res.Error != nil {
		c.AbortWithError(http.StatusInternalServerError, res.Error)
	}

	return res
}