package user

import (
	"context"
	"faceit-backend-test/internal/apierr"
	"faceit-backend-test/internal/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const route = "/users"

type Service interface {
	Create(ctx context.Context, request CreateUserRequest) (CreateUserResponse, error)
	Update(ctx context.Context, request UpdateUserRequest) (UpdateUserResponse, error)
	DeleteById(ctx context.Context, request DeleteUserByIdRequest) (DeleteUserResponse, error)
	GetMany(ctx context.Context) (GetUsersManyResponse, error)
}

type controller struct {
	service Service
}

var _ router.Controller = (*controller)(nil)

type ControllerOpts func(*controller)

func NewController(opts ...ControllerOpts) router.Controller {
	c := &controller{}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithService(service Service) ControllerOpts {
	return func(controller *controller) {
		controller.service = service
	}
}

func (c *controller) Register(r *gin.RouterGroup) {
	r.GET(route, c.getUsersMany)
	r.POST(route, c.createUser)
	r.PATCH(fmt.Sprintf("%v/:id", route), c.updateUser)
	r.DELETE(fmt.Sprintf("%v/:id", route), c.deleteUserById)
}

func (c *controller) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.decodeError(ctx, apierr.BadRequest(err.Error()))
		return
	}

	resp, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		c.decodeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

func (c *controller) updateUser(ctx *gin.Context) {
	var req UpdateUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.decodeError(ctx, err)
		return
	}

	req.Id = ctx.Param("id")
	resp, err := c.service.Update(ctx.Request.Context(), req)
	if err != nil {
		c.decodeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *controller) deleteUserById(ctx *gin.Context) {
	var req DeleteUserByIdRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		c.decodeError(ctx, err)
		return
	}

	resp, err := c.service.DeleteById(ctx.Request.Context(), req)
	if err != nil {
		c.decodeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *controller) getUsersMany(ctx *gin.Context) {
	resp, err := c.service.GetMany(ctx)
	if err != nil {
		c.decodeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *controller) decodeError(ctx *gin.Context, err error) {
	apiError, ok := err.(apierr.ApiError)
	if !ok {
		apiError = apierr.InternalServerError()
	}

	ctx.JSON(apiError.StatusCode, apiError)
}
