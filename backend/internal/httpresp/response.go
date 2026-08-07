// Package httpresp provides a consistent JSON envelope for API responses.
package httpresp

import "github.com/gin-gonic/gin"

func OK(c *gin.Context, code int, data any) {
	c.JSON(code, gin.H{"data": data})
}

func Err(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}
