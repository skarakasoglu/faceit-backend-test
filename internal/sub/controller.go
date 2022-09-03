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

// controller operations related to subscription to
// events in the app are handled in this controller
// @tag.name SubscribeController
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

// Register registers the endpoints to the given router group
func (c *controller) Register(r *gin.RouterGroup) {
	r.POST(route, c.subscribe)
}

// subscribe godoc
// @Summary creates a subscribe request to given topic and sends it to verification queue, returns the subscription details.
// @tags SubscribeController
// @Accept json
// @Produce json
// @Param SubscribeRequest body SubscribeRequest true "subscription details"
// @Success 200 {object} SubscribeResponse
// @Failure 400 {object} apierr.ApiError
// @Failure 500 {object} apierr.ApiError
// @Router /v1/subscribe [post]
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
