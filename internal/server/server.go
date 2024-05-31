package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"html/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/madswillem/recipeApp_Backend_Go/internal/database"
	"github.com/madswillem/recipeApp_Backend_Go/internal/middleware"
	"gorm.io/gorm"
)

type Server struct {
	port int
	DB *gorm.DB
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewSr := &Server{
		port: port,
		DB: database.ConnectToDB(),
	}

	// Declare Sr config
	sr := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewSr.port),
		Handler:      NewSr.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return sr
}
func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	tmpl := template.Must(template.New("main").ParseGlob("web/templates/**/*"))
	r.SetHTMLTemplate(tmpl)
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404/index", gin.H{
			"pageTitle": "404 Page not found",
		})
	})
	r.Use(middleware.CORSMiddleware())

	r.POST("/create", s.AddRecipe)
	r.GET("/get", s.GetAll)
	r.GET("/getbyid/:id", s.GetById)
	r.PATCH("/update/:id", s.UpdateRecipe)
	r.DELETE("/delete/:id", s.DeleteRecipe)
	r.POST("/filter", s.Filter)
	r.GET("/select/:id", middleware.User, s.Select)
	r.GET("/deselect/:id", middleware.User, s.Deselect)
	r.GET("/colormode/:type", s.Colormode)

	r.PATCH("/account/update", middleware.User, s.UpadateUser)

	r.GET("/", s.RenderHome)
	r.GET("/account", s.RenderAcount)
	r.GET("/tutorials", s.RenderTutorial)
	r.GET("/recipe/:id", s.RenderProductpage)

	r.GET("/get/home", s.GetHome)
	r.GET("/get/account", s.GetAccount)
	r.GET("/get/recipe/:id", s.GetRecipe)

	r.GET("/style/:filename", s.GetStyles)
	r.GET("/imgs/:filename", s.GetImgs)
	r.GET("/scripts/:filename", s.GetScripts)

	return r
}

