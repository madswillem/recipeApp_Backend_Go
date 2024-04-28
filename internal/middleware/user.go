package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
)

func User(c *gin.Context) {
	var user models.UserModel
	cookie, err := c.Cookie("user")
	if err != nil || cookie == "" {
		user.Create("unkwon")
		println("WHY")
		c.SetCookie("user", user.Cookie, 31536000, "/", "localhost:8080", false, true)
	} else {
		println(cookie)
		user.Cookie = cookie
		err := user.GetByCookie()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		}
	}
	c.Set("user", user)
	c.Next()
}
