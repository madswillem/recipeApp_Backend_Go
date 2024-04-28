package main

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madswillem/recipeApp_Backend_Go/internal/controllers"
	"github.com/madswillem/recipeApp_Backend_Go/internal/initializers"
	"github.com/madswillem/recipeApp_Backend_Go/internal/middleware"
	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

	initializers.DB.AutoMigrate(&models.RecipeSchema{})
	initializers.DB.AutoMigrate(&models.IngredientsSchema{})
	initializers.DB.AutoMigrate(&models.RatingStruct{})
	initializers.DB.AutoMigrate(&models.NutritionalValue{})
	initializers.DB.AutoMigrate(&models.DietSchema{})
	initializers.DB.AutoMigrate(&models.IngredientDBSchema{})
}

func main() {
	r := gin.Default()
	tmpl := template.Must(template.New("main").ParseGlob("web/templates/**/*"))
	r.SetHTMLTemplate(tmpl)
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404/index", gin.H{
			"pageTitle": "404 Page not found",
		})
	})
	r.Use(middleware.CORSMiddleware())

	r.POST("/create", controllers.AddRecipe)
	r.GET("/get", controllers.GetAll)
	r.GET("/getbyid/:id", controllers.GetById)
	r.PATCH("/update/:id", controllers.UpdateRecipe)
	r.DELETE("/delete/:id", controllers.DeleteRecipe)
	r.POST("/filter", controllers.Filter)
	r.GET("/select/:id", controllers.Select)
	r.GET("/deselect/:id", controllers.Deselect)
	r.GET("/colormode/:type", controllers.Colormode)

	r.Run()
}
