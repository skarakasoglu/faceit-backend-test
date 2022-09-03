package health

import (
	"faceit-backend-test/internal/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

const route = "/health"

// controller it handles the health check
// @tag.name HealthController
type controller struct{}

func NewController() router.Controller {
	return &controller{}
}

func (c *controller) Register(r *gin.RouterGroup) {
	r.GET(route, c.healthCheck)
}

// healthCheck godoc
// @Summary checks the status of the service
// @tags HealthController
// @Produce json
// @Success 200 {object} Response
// @Router /v1/health [get]
func (c *controller) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{Status: true})
}
