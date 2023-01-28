package main

import (
	"github.com/gin-gonic/gin"
	"rezeptapp.ml/goApp/controllers"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/middleware"
)

func init()  {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main()  {
	r := gin.Default()
	r.POST("/create", middleware.CORSMiddleware(), controllers.AddRecipe)
	r.GET("/get", middleware.CORSMiddleware(), controllers.GetAll)
	r.GET("/getbyid/:id", middleware.CORSMiddleware(), controllers.GetById)
	r.GET("/select/:id", middleware.CORSMiddleware(), controllers.Select)
	r.GET("/deselect/:id", middleware.CORSMiddleware(), controllers.Deselect)
	r.GET("/colormode/:type", middleware.CORSMiddleware(), controllers.Colormode)
	r.Run()
}