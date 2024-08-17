package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/madswillem/recipeApp_Backend_Go/internal/database"
	"gorm.io/gorm"
)

const MethodGET = "GET"
const MethodPost = "POST"

type InnitFuncs func(*Server) error
type ExtraControllers struct {
	Function   func(*gin.Context)
	Middleware func(*gin.Context)
	Route      string
	Method     string
}

type Config struct {
	Innit       []InnitFuncs
	Controllers []ExtraControllers
	DBConf      gorm.Config
}

type Server struct {
	port   int
	DB     *gorm.DB
	NewDB  *sqlx.DB
	config *Config
}

func NewServer(config *Config) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:   port,
		DB:     database.ConnectToGORMDB(&config.DBConf),
		NewDB:  database.ConnectToDB(&sqlx.Conn{}),
		config: config,
	}

	for _, fnc := range NewServer.config.Innit {
		err := fnc(NewServer)
		fmt.Println(err)
	}

	if NewServer.DB == nil {
		panic("db is nil")
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
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
	r.Use(s.CORSMiddleware())

	for _, controller := range s.config.Controllers {
		if controller.Route == MethodGET {
			r.GET(controller.Route, controller.Function)
		}
		if controller.Method == MethodPost {
			r.GET(controller.Route, controller.Function)
		}
	}

	r.POST("/create", s.AddRecipe)
	r.POST("/create_ingredient", s.AddIngredient)
	r.GET("/get", s.GetAll)
	r.GET("/getbyid/:id", s.GetById)
	r.PATCH("/update/:id", s.UpdateRecipe)
	r.DELETE("/delete/:id", s.DeleteRecipe)
	r.POST("/filter", s.Filter)
	r.GET("/select/:id", s.UserMiddleware, s.Select)
	r.GET("/deselect/:id", s.UserMiddleware, s.Deselect)
	r.GET("/colormode/:type", s.Colormode)

	r.PATCH("/account/update", s.UserMiddleware, s.UpadateUser)
	r.GET("/creategroup", s.UserMiddleware, s.CreateGroup)
	r.GET("/recommendation", s.UserMiddleware, s.GetRecommendation)

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
