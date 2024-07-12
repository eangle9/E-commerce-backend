package middleware

import (
	"Eccomerce-website/internal/core/entity"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(middlewareLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, exist := c.Get("requestID")
		if !exist {
			err := errors.New("requestID not found in gin context")
			errorResponse := entity.AppInternalError.Wrap(err, "failed to get requestID")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errorResponse.Error()})
			middlewareLogger.Error("requestID error",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "middlewareLayer"),
				zap.String("function", "LoggerMiddleware"),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return
		}

		requestID, ok := id.(string)
		if !ok {
			err := errors.New("unable to convert type any to string")
			errorResponse := entity.AppInternalError.Wrap(err, "failed to convert requestId type any to string")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errorResponse.Error()})
			middlewareLogger.Error("requestId is not exist in the context",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "middlewareLayer"),
				zap.String("function", "LoggerMiddleware"),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return
		}

		c.Next()

		middlewareLogger.Info("request processed",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("requestID", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("statusCode", c.Writer.Status()),
			zap.String("clientIP", c.ClientIP()),
			zap.String("userAgent", c.Request.UserAgent()),
		)

	}
}
