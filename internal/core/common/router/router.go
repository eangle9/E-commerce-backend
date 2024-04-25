package router

import "github.com/gin-gonic/gin"

type Router struct {
	*gin.Engine
}

func NewRouter(instance *gin.Engine) *Router {
	return &Router{
		instance,
	}
}
