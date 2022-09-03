package router

import (
	_ "faceit-backend-test/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Controller is an interface for routes that will be served in the router
type Controller interface {
	Register(group *gin.RouterGroup)
}

// NewHTTPRouter creates a new router and registers all the routes
func NewHTTPRouter(routes ...Controller) *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
