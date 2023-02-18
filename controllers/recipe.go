package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/middleware"
	"rezeptapp.ml/goApp/models"
)

func GetAll(c *gin.Context) {
	var recipes []models.RecipeSchema

	result := initializers.DB.Preload("Ingredients").Find(&recipes)

	if result.Error != nil {
		panic(result.Error)
	}

	c.JSON(http.StatusOK, recipes)
}

func AddRecipe(c *gin.Context) {
	var body models.RecipeSchema

	err := c.Bind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Failed to read body",
			"errMessage": err.Error(),
		})
		return
	}

	result := initializers.DB.Create(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Failed to read body",
			"errMessage": err.Error(),
			"result":     result,
		})
		return
	}

	c.JSON(http.StatusCreated, body)
}

func GetById(c *gin.Context) {
	result := middleware.GetDataByID(c.Param("id"), c)

	c.JSON(http.StatusOK, result)
}

func Select(c *gin.Context) {
	res := middleware.UpdateSelected(c.Param("id"), +1, c)

	c.JSON(http.StatusOK, res)
}

func Deselect(c *gin.Context) {
	res := middleware.UpdateSelected(c.Param("id"), -1, c)

	c.JSON(http.StatusOK, res)
}

func Colormode(c *gin.Context) {
	if c.Param("type") == "get" {
		cookie, err := c.Cookie("type")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"type": cookie,
		})
	} else if c.Param("type") == "dark" {
		c.SetCookie("type", "dark", 999999999999999999, "/", "localhost", false, true)
		c.Status(http.StatusAccepted)
	} else if c.Param("type") == "light" {
		c.SetCookie("type", "light", 999999999999999999, "/", "localhost", false, true)
		c.Status(http.StatusAccepted)
	} else {
		c.Status(http.StatusBadRequest)
	}
}
