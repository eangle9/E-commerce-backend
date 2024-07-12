package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := fmt.Sprintf("%d", time.Now().UnixNano())
		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}
