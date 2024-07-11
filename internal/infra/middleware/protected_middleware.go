package middleware

import (
	jwttoken "Eccomerce-website/internal/core/common/jwt_token"
	"Eccomerce-website/internal/core/entity"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ProtectedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		if bearerToken == "" {
			err := errors.New("missing authorization header")
			errorResponse := entity.InvalidCredentials.Wrap(err, "token is empty in the authorization header").WithProperty(entity.StatusCode, 401)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errorResponse.Error()})
			return
		}

		token := strings.Split(bearerToken, " ")[1]
		id, role, err := jwttoken.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("userId", id)
		c.Set("role", role)

		c.Next()
	}
}
