package middleware

import (
	jwttoken "Eccomerce-website/internal/core/common/jwt_token"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ProtectedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		if bearerToken == "" {
			errorResponse := response.Response{
				Status:       http.StatusUnauthorized,
				ErrorType:    errorcode.Unauthorized,
				ErrorMessage: "missing authorization header",
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse)
			return
		}

		token := strings.Split(bearerToken, " ")[1]
		id, role, err := jwttoken.VerifyToken(token)
		if err != nil {
			errorResponse := response.Response{
				Status:       http.StatusUnauthorized,
				ErrorType:    errorcode.Unauthorized,
				ErrorMessage: err.Error(),
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse)
			return
		}

		c.Set("userId", id)
		c.Set("role", role)

		c.Next()
	}
}
