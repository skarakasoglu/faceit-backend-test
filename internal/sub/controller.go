package sub

import (
	"context"
	"faceit-backend-test/internal/apierr"
	"faceit-backend-test/internal/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

const route = "/subscribe"

type Service interface {
	Subscribe(ctx context.Context, request SubscribeRequest) (SubscribeResponse, error)
}

type controller struct {
	service Service
}

type ControllerOpts func(*controller)

func NewController(opts ...ControllerOpts) router.Controller {
	c := &controller{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithService(service Service) ControllerOpts {
	return func(c *controller) {
		c.service = service
	}
}

func (c *controller) Register(r *gin.RouterGroup) {
	r.POST(route, c.subscribe)
}

func (c *controller) subscribe(ctx *gin.Context) {
	var req SubscribeRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.decodeError(ctx, apierr.BadRequest(err.Error()))
		return
	}

	resp, err := c.service.Subscribe(ctx.Request.Context(), req)
	if err != nil {
		c.decodeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

func (c *controller) decodeError(ctx *gin.Context, err error) {
	apiErr, ok := err.(apierr.ApiError)
	if !ok {
		apiErr = apierr.InternalServerError()
	}

	ctx.JSON(apiErr.StatusCode, apiErr)
}
