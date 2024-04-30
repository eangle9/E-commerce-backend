package middleware

import (
	"Eccomerce-website/internal/core/model/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err, exist := c.Get("error")
		if exist {
			var statuscode int
			errorResponse := err.(response.Response)

			switch errorResponse.ErrorType {
			case "NOT_FOUND_ERROR":
				statuscode = http.StatusNotFound
			case "VALIDATION_ERROR":
				statuscode = http.StatusBadRequest
			case "INVALID_REQUEST":
				statuscode = http.StatusBadRequest
			case "INTERNAL_SERVER_ERROR":
				statuscode = http.StatusInternalServerError
			case "UNAUTHORIZED":
				statuscode = http.StatusUnauthorized
			default:
				statuscode = errorResponse.Status

			}

			c.Header("Content-Type", "application/json")
			c.AbortWithStatusJSON(statuscode, errorResponse)
		}
	}
}
