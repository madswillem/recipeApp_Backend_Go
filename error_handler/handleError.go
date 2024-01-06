package error_handler

import "github.com/gin-gonic/gin"
func HandleError(c *gin.Context, statusCode int, errorMessage string, err error) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error":      errorMessage,
		"errMessage": err.Error(),
	})
}