package middleware

import (
	"Eccomerce-website/internal/core/entity"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TimeoutMiddleware(timeoutLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		done := make(chan struct{})

		go func() {
			defer close(done)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		}()

		id, exist := c.Get("requestID")
		if !exist {
			err := errors.New("requestID not found in gin context")
			errorResponse := entity.AppInternalError.Wrap(err, "failed to get requestID")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errorResponse.Error()})
			timeoutLogger.Error("requestID error",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "middlewareLayer"),
				zap.String("function", "ProtectedMiddleware"),
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
			timeoutLogger.Error("requestId is not exist in the context",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "middlewareLayer"),
				zap.String("function", "ProtectedMiddleware"),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return
		}

		select {
		case <-done:
			return
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				errorResponse := entity.TimeoutError.Wrap(ctx.Err(), "request timeout")
				c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{"error": errorResponse.Error()})
				timeoutLogger.Error("timeout error",
					zap.String("timestamp", time.Now().Format(time.RFC3339)),
					zap.String("layer", "middlewareLayer"),
					zap.String("function", "TimeoutMiddleware"),
					zap.String("requestID", requestID),
					zap.Error(errorResponse),
					zap.Stack("stacktrace"),
				)
				return
			}
		}
	}
}
