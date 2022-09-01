package health

import (
	"faceit-backend-test/internal/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

const route = "/health"

type controller struct{}

func NewController() router.Controller {
	return &controller{}
}

func (c *controller) Register(r *gin.RouterGroup) {
	r.GET(route, c.healthCheck)
}

func (c *controller) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{Status: true})
}
