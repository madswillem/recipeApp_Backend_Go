package main

import (
	"github.com/gin-gonic/gin"
	"rezeptapp.ml/goApp/controllers"
	"rezeptapp.ml/goApp/initializers"
)

func init()  {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main()  {
	r := gin.Default()
	r.POST("/create", controllers.AddRecipe)
	r.GET("/get", controllers.GetAll)
	r.GET("/getbyid/:id", controllers.GetById)
	r.GET("/select/:id", controllers.Select)
	r.GET("/deselect/:id", controllers.Deselect)
	r.GET("/colormode/:type", controllers.Colormode)
	r.Run()
}