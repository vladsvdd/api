package handler

import (
	"github.com/gin-gonic/gin"
)

func newErrorResponse(c *gin.Context, statusCode, code int, message string) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"code":    code,
		"message": message,
		"details": "{}",
	})
}
