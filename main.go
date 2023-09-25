package main

import (
	"net/http"
	"html/template"

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
	tmpl := template.Must(template.New("main").ParseGlob("templates/**/*"))
	r.SetHTMLTemplate(tmpl)
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"pageTitle": "404 Page not found",
		})
	})

	r.POST("/create", middleware.CORSMiddleware(), controllers.AddRecipe)
	r.GET("/recommended", middleware.CORSMiddleware(), controllers.Recomend)
	r.GET("/get", middleware.CORSMiddleware(), controllers.GetAll)
	r.GET("/getbyid/:id", middleware.CORSMiddleware(), controllers.GetById)
	r.POST("/filter", middleware.CORSMiddleware(), controllers.Filter)
	r.GET("/select/:id", middleware.CORSMiddleware(), controllers.Select)
	r.GET("/deselect/:id", middleware.CORSMiddleware(), controllers.Deselect)
	r.GET("/colormode/:type", middleware.CORSMiddleware(), controllers.Colormode)

	r.GET("/", controllers.RenderHome)
	r.GET("/account", controllers.RenderAcount)
	r.GET("/recipe/:id", controllers.RenderProductpage)

	r.GET("/get/home", controllers.GetHome)
	r.GET("/get/account", controllers.GetAccount)
	r.GET("/get/recipe/:id", controllers.GetRecipe)

	r.GET("/style/:filename", controllers.GetStyles)
	r.GET("/imgs/:filename", controllers.GetImgs)
	
	r.Run()
}
