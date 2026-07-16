package handlers

import (
	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func RespondError(c *gin.Context, err *AppError) {
	c.JSON(err.Code, gin.H{"error": err.Message})
}

func RespondJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}
