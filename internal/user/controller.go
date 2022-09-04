package user

import (
	"context"
	"encoding/json"
	_ "faceit-backend-test/docs"
	"faceit-backend-test/internal/apierr"
	"faceit-backend-test/internal/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
)

const (
	defaultPage    = 1
	defaultPerPage = 10
)
const route = "/users"

type Service interface {
	Create(ctx context.Context, request CreateUserRequest) (CreateUserResponse, error)
	Update(ctx context.Context, request UpdateUserRequest) (UpdateUserResponse, error)
	DeleteById(ctx context.Context, request DeleteUserByIdRequest) (DeleteUserResponse, error)
	GetMany(ctx context.Context, request GetUsersManyRequest) (GetUsersManyResponse, error)
}

// controller it handles the operations related to users
// @tag.name UserController
type controller struct {
	service Service
}

var _ router.Controller = (*controller)(nil)

type ControllerOpts func(*controller)

// NewController create new instance of user.Controller with options
func NewController(opts ...ControllerOpts) *controller {
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

// Register it registers the routes and handlers
// to the router group passed as an argument.
func (c *controller) Register(r *gin.RouterGroup) {
	r.GET(route, c.GetUsersMany)
	r.POST(route, c.CreateUser)
	r.PATCH(fmt.Sprintf("%v/:id", route), c.UpdateUser)
	r.DELETE(fmt.Sprintf("%v/:id", route), c.DeleteUserById)
}

// CreateUser godoc
// @Summary creates a user
// @tags UserController
// @Accept json
// @Produce json
// @Param CreateUserRequest body CreateUserRequest true "user details"
// @Success 200 {object} CreateUserResponse
// @Failure 400 {object} apierr.ApiError
// @Failure 500 {object} apierr.ApiError
// @Router /v1/users [post]
func (c *controller) CreateUser(ctx *gin.Context) {
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

// UpdateUser godoc
// @Summary updates a user having id provided in path param
// @tags UserController
// @Accept json
// @Produce json
// @Param id path string true "id of the user"
// @Param UpdateUserRequest body UpdateUserRequest true "user details"
// @Success 200 {object} UpdateUserResponse
// @Failure 400 {object} apierr.ApiError
// @Failure 500 {object} apierr.ApiError
// @Router /v1/users/{id} [patch]
func (c *controller) UpdateUser(ctx *gin.Context) {
	var req UpdateUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.decodeError(ctx, apierr.BadRequest(err.Error()))
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

// DeleteUserById godoc
// @Summary deletes the user having id provided in path param
// @tags UserController
// @Accept json
// @Produce json
// @Param id path string true "id of the user"
// @Success 200 {object} DeleteUserResponse
// @Failure 400 {object} apierr.ApiError
// @Failure 500 {object} apierr.ApiError
// @Router /v1/users/{id} [delete]
func (c *controller) DeleteUserById(ctx *gin.Context) {
	var req DeleteUserByIdRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		c.decodeError(ctx, apierr.BadRequest(err.Error()))
		return
	}

	resp, err := c.service.DeleteById(ctx.Request.Context(), req)
	if err != nil {
		c.decodeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetUsersMany godoc
// @Summary returns the users with respect to the pagination and filter parameters
// @tags UserController
// @Accept json
// @Produce json
// @Param page query int false "page number that will be returned" Default(1)
// @Param perPage query int false "how many rows are returned by page" Default(10)
// @Param filter query string false "filtering parameters that will be used while fetching the users" example({"country": "UK", "first_name": "Alisson"})
// @Success 200 {object} GetUsersManyResponse
// @Failure 400 {object} apierr.ApiError
// @Failure 500 {object} apierr.ApiError
// @Router /v1/users [get]
func (c *controller) GetUsersMany(ctx *gin.Context) {
	var req GetUsersManyRequest
	var err error
	req.Page = defaultPage
	req.PerPage = defaultPerPage

	if page := ctx.Query("page"); page != "" {
		req.Page, err = strconv.Atoi(page)
		if err != nil {
			c.decodeError(ctx, apierr.BadRequest(err.Error()))
			return
		}
	}

	if perPage := ctx.Query("perPage"); perPage != "" {
		req.PerPage, err = strconv.Atoi(perPage)
		if err != nil {
			c.decodeError(ctx, apierr.BadRequest(err.Error()))
			return
		}
	}

	if filter := ctx.Query("filter"); filter != "" {
		decodedFilter, err := url.QueryUnescape(filter)
		if err != nil {
			c.decodeError(ctx, apierr.BadRequest(err.Error()))
			return
		}

		err = json.Unmarshal([]byte(decodedFilter), &req.Filter)
		if err != nil {
			c.decodeError(ctx, apierr.BadRequest(err.Error()))
			return
		}
	}

	resp, err := c.service.GetMany(ctx, req)
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
