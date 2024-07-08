package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
)

func (s *Server) UserMiddleware(c *gin.Context) {
	var user models.UserModel
	cookie, err := c.Cookie("user")
	if err != nil || cookie == "" {
		user.Create(s.NewDB ,"unkwon")
		c.SetCookie("user", user.Cookie, 31536000, "/", "localhost", false, true)
	} else {
		user.Cookie = cookie
		err := user.GetByCookie(s.DB)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		}
	}
	c.Set("user", user)
	c.Next()
}
