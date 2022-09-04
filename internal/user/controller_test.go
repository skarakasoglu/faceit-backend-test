package user

import (
	"bytes"
	"context"
	"encoding/json"
	"faceit-backend-test/internal/apierr"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type serviceMock struct {
	createMock     func(context.Context, CreateUserRequest) (CreateUserResponse, error)
	updateMock     func(context.Context, UpdateUserRequest) (UpdateUserResponse, error)
	deleteByIdMock func(context.Context, DeleteUserByIdRequest) (DeleteUserResponse, error)
	getManyMock    func(context.Context, GetUsersManyRequest) (GetUsersManyResponse, error)
}

func (s *serviceMock) Create(ctx context.Context, request CreateUserRequest) (CreateUserResponse, error) {
	return s.createMock(ctx, request)
}

func (s *serviceMock) Update(ctx context.Context, request UpdateUserRequest) (UpdateUserResponse, error) {
	return s.updateMock(ctx, request)
}

func (s *serviceMock) DeleteById(ctx context.Context, request DeleteUserByIdRequest) (DeleteUserResponse, error) {
	return s.deleteByIdMock(ctx, request)
}

func (s *serviceMock) GetMany(ctx context.Context, request GetUsersManyRequest) (GetUsersManyResponse, error) {
	return s.getManyMock(ctx, request)
}

func TestController_CreateUser(t *testing.T) {
	mockService := &serviceMock{}
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		expected := CreateUserResponse{User{
			Id:        e.Id,
			FirstName: e.FirstName,
			LastName:  e.LastName,
			Nickname:  e.Nickname,
			Password:  e.Password,
			Email:     e.Email,
			Country:   e.Country,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		}}
		req := CreateUserRequest{
			FirstName: e.FirstName,
			LastName:  e.LastName,
			Nickname:  e.Nickname,
			Password:  e.Password,
			Email:     e.Email,
			Country:   e.Country,
		}

		mockService.createMock = func(ctx context.Context, request CreateUserRequest) (CreateUserResponse, error) {
			assert.EqualValues(t, req, request)

			return expected, nil
		}

		controller := NewController(WithService(mockService))

		rr := httptest.NewRecorder()

		router := gin.Default()
		router.POST(route, controller.CreateUser)

		reqBody, err := json.Marshal(req)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(expected)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
	})

	t.Run("should return bad request when a required field is missing", func(t *testing.T) {
		req := CreateUserRequest{
			FirstName: e.FirstName,
			Nickname:  e.Nickname,
			Password:  e.Password,
			Email:     e.Email,
			Country:   e.Country,
		}
		expected := apierr.BadRequest("invalid request body")

		controller := NewController(WithService(mockService))

		rr := httptest.NewRecorder()

		router := gin.Default()
		router.POST(route, controller.CreateUser)

		reqBody, err := json.Marshal(req)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		var actualResp apierr.ApiError
		err = json.Unmarshal(rr.Body.Bytes(), &actualResp)

		assert.Equal(t, expected.StatusCode, rr.Code)
		assert.Equal(t, expected.Code, actualResp.Code)
	})

	t.Run("should return internal error when service fails", func(t *testing.T) {
		req := CreateUserRequest{
			FirstName: e.FirstName,
			Nickname:  e.Nickname,
			Password:  e.Password,
			LastName:  e.LastName,
			Email:     e.Email,
			Country:   e.Country,
		}

		expected := repositoryError(fmt.Errorf("mock error"))
		mockService.createMock = func(ctx context.Context, request CreateUserRequest) (CreateUserResponse, error) {
			return CreateUserResponse{}, expected
		}
		controller := NewController(WithService(mockService))

		rr := httptest.NewRecorder()

		router := gin.Default()
		router.POST(route, controller.CreateUser)

		reqBody, err := json.Marshal(req)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		var actualResp apierr.ApiError
		err = json.Unmarshal(rr.Body.Bytes(), &actualResp)

		assert.Equal(t, expected.StatusCode, rr.Code)
		assert.Equal(t, expected.Code, actualResp.Code)
		assert.Equal(t, expected.Message, actualResp.Message)
	})
}
