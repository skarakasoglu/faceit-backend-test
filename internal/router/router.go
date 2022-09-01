package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller interface {
	Register(group *gin.RouterGroup)
}

func NewHTTPRouter(routes ...Controller) http.Handler {
	router := gin.Default()

	createApiV1(router, routes...)
	return router
}

func createApiV1(router *gin.Engine, routes ...Controller) *gin.RouterGroup {
	v1 := router.Group("v1")
	for _, route := range routes {
		route.Register(v1)
	}

	return v1
}
