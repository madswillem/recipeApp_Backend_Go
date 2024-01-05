package tools

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"rezeptapp.ml/goApp/models"
)

func GetSimilar(c *gin.Context, recipes []models.RecipeSchema) []models.RecipeSchema {
	var list []models.RecipeSchema
	fmt.Println(list)
	return recipes
}