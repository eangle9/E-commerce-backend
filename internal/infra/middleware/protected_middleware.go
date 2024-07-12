package middleware

import (
	jwttoken "Eccomerce-website/internal/core/common/jwt_token"
	"Eccomerce-website/internal/core/entity"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ProtectedMiddleware(middlewareLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		id, exist := c.Get("requestID")
		if !exist {
			err := errors.New("requestID not found in gin context")
			errorResponse := entity.AppInternalError.Wrap(err, "failed to get requestID")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errorResponse.Error()})
			middlewareLogger.Error("requestID error",
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
			middlewareLogger.Error("requestId is not exist in the context",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "middlewareLayer"),
				zap.String("function", "ProtectedMiddleware"),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return
		}

		bearerToken := c.Request.Header.Get("Authorization")
		if bearerToken == "" {
			err := errors.New("missing authorization header")
			errorResponse := entity.InvalidCredentials.Wrap(err, "token is empty in the authorization header").WithProperty(entity.StatusCode, 401)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errorResponse.Error()})
			middlewareLogger.Error("bearerToken is empty",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "middlewareLayer"),
				zap.String("function", "ProtectedMiddleware"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return
		}

		token := strings.Split(bearerToken, " ")[1]
		id, role, err := jwttoken.VerifyToken(ctx, token, middlewareLogger, requestID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("userId", id)
		c.Set("role", role)

		c.Next()
	}
}
