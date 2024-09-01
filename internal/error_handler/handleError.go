package error_handler

import "github.com/gin-gonic/gin"

type APIError struct {
	Message string
	Code    int
	Errors  []error
}
type Errors struct {
	Errors []error
}

func New(message string, code int, errors error) *APIError {
	return &APIError{
		Message: message,
		Code:    code,
		Errors:  []error{errors},
	}
}

func HandleError(c *gin.Context, statusCode int, errorMessage string, err []error) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"errors":     err[0].Error(),
		"errMessage": errorMessage,
	})
}
