package middleware

import (
	"Eccomerce-website/internals/core/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		errors := c.Errors

		if len(errors) > 0 {
			for _, e := range errors {
				if err, ok := e.Err.(*errorx.Error); ok {
					statuscode, exist := err.Property(entity.StatusCode)
					if !exist {
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}

					status, ok := statuscode.(int)
					if !ok {
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "status code is not an integer"})
						return
					}

					c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
					return
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
			}
		}
	}
}
