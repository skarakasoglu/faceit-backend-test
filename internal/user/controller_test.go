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
	"strconv"
	"testing"
)

type mockService struct {
	createMock     func(context.Context, CreateUserRequest) (CreateUserResponse, error)
	updateMock     func(context.Context, UpdateUserRequest) (UpdateUserResponse, error)
	deleteByIdMock func(context.Context, DeleteUserByIdRequest) (DeleteUserResponse, error)
	getManyMock    func(context.Context, GetUsersManyRequest) (GetUsersManyResponse, error)
}

func (s *mockService) Create(ctx context.Context, request CreateUserRequest) (CreateUserResponse, error) {
	return s.createMock(ctx, request)
}

func (s *mockService) Update(ctx context.Context, request UpdateUserRequest) (UpdateUserResponse, error) {
	return s.updateMock(ctx, request)
}

func (s *mockService) DeleteById(ctx context.Context, request DeleteUserByIdRequest) (DeleteUserResponse, error) {
	return s.deleteByIdMock(ctx, request)
}

func (s *mockService) GetMany(ctx context.Context, request GetUsersManyRequest) (GetUsersManyResponse, error) {
	return s.getManyMock(ctx, request)
}

func TestController_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	controller := NewController()

	t.Run("", func(t *testing.T) {
		controller.Register(&router.RouterGroup)

		rr := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, "/users", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.NotEqual(t, http.StatusNotFound, rr.Code)
	})
}

func TestController_CreateUser(t *testing.T) {
	mockService := &mockService{}
	gin.SetMode(gin.TestMode)
	controller := NewController(WithService(mockService))
	router := gin.Default()
	router.POST(route, controller.CreateUser)

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

		reqBody, err := json.Marshal(req)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
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

		reqBody, err := json.Marshal(req)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
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

		reqBody, err := json.Marshal(req)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, request)

		var actualResp apierr.ApiError
		err = json.Unmarshal(rr.Body.Bytes(), &actualResp)

		assert.Equal(t, expected.StatusCode, rr.Code)
		assert.Equal(t, expected.Code, actualResp.Code)
		assert.Equal(t, expected.Message, actualResp.Message)
	})
}

func TestController_UpdateUser(t *testing.T) {
	mockService := &mockService{}
	controller := NewController(WithService(mockService))

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PATCH(fmt.Sprintf("%v/:id", route), controller.UpdateUser)

	t.Run("success", func(t *testing.T) {
		expected := UpdateUserResponse{User{
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
		req := UpdateUserRequest{
			Id:        e.Id,
			FirstName: e.FirstName,
			LastName:  e.LastName,
			Nickname:  e.Nickname,
			Password:  e.Password,
			Email:     e.Email,
			Country:   e.Country,
		}

		mockService.updateMock = func(ctx context.Context, request UpdateUserRequest) (UpdateUserResponse, error) {
			assert.EqualValues(t, req, request)

			return expected, nil
		}

		reqBody, err := json.Marshal(req)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%v", e.Id), bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(expected)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
	})

	t.Run("should return bad request when a required field is missing", func(t *testing.T) {
		req := UpdateUserRequest{
			FirstName: e.FirstName,
			Nickname:  e.Nickname,
			Password:  e.Password,
			Email:     e.Email,
			Country:   e.Country,
		}
		expected := apierr.BadRequest("invalid request body")

		reqBody, err := json.Marshal(req)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%v", e.Id), bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, request)

		var actualResp apierr.ApiError
		err = json.Unmarshal(rr.Body.Bytes(), &actualResp)

		assert.Equal(t, expected.StatusCode, rr.Code)
		assert.Equal(t, expected.Code, actualResp.Code)
	})

	t.Run("should return internal error when service fails", func(t *testing.T) {
		req := UpdateUserRequest{
			FirstName: e.FirstName,
			Nickname:  e.Nickname,
			Password:  e.Password,
			LastName:  e.LastName,
			Email:     e.Email,
			Country:   e.Country,
		}

		expected := repositoryError(fmt.Errorf("mock error"))
		mockService.updateMock = func(ctx context.Context, request UpdateUserRequest) (UpdateUserResponse, error) {
			return UpdateUserResponse{}, expected
		}

		reqBody, err := json.Marshal(req)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%v", e.Id), bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, request)

		var actualResp apierr.ApiError
		err = json.Unmarshal(rr.Body.Bytes(), &actualResp)

		assert.Equal(t, expected.StatusCode, rr.Code)
		assert.Equal(t, expected.Code, actualResp.Code)
		assert.Equal(t, expected.Message, actualResp.Message)
	})
}

func TestController_DeleteUserById(t *testing.T) {
	mockService := &mockService{}
	controller := NewController(WithService(mockService))

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE(fmt.Sprintf("%v/:id", route), controller.DeleteUserById)

	t.Run("success", func(t *testing.T) {
		expected := DeleteUserResponse{
			Id: e.Id,
		}
		req := DeleteUserByIdRequest{
			Id: e.Id,
		}

		mockService.deleteByIdMock = func(ctx context.Context, request DeleteUserByIdRequest) (DeleteUserResponse, error) {
			assert.EqualValues(t, req, request)

			return expected, nil
		}

		reqBody, err := json.Marshal(req)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%v", e.Id), bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(expected)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
	})

	t.Run("should return internal error when service fails", func(t *testing.T) {
		req := DeleteUserByIdRequest{
			Id: e.Id,
		}

		expected := repositoryError(fmt.Errorf("mock error"))
		mockService.deleteByIdMock = func(ctx context.Context, request DeleteUserByIdRequest) (DeleteUserResponse, error) {
			return DeleteUserResponse{}, expected
		}

		reqBody, err := json.Marshal(req)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%v", e.Id), bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, request)

		var actualResp apierr.ApiError
		err = json.Unmarshal(rr.Body.Bytes(), &actualResp)

		assert.Equal(t, expected.StatusCode, rr.Code)
		assert.Equal(t, expected.Code, actualResp.Code)
		assert.Equal(t, expected.Message, actualResp.Message)
	})
}

func TestController_GetUsersMany(t *testing.T) {
	mockService := &mockService{}
	controller := NewController(WithService(mockService))

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET(fmt.Sprintf("%v", route), controller.GetUsersMany)

	t.Run("success with default query params", func(t *testing.T) {
		req := GetUsersManyRequest{
			Page:    1,
			PerPage: 10,
			Filter:  User{},
		}

		expected := GetUsersManyResponse{}
		expected.Users = make([]User, len(entities))
		for i, val := range entities {
			expected.Users[i] = User{
				Id:        val.Id,
				FirstName: val.FirstName,
				LastName:  val.LastName,
				Nickname:  val.Nickname,
				Password:  val.Password,
				Email:     val.Email,
				Country:   val.Country,
				CreatedAt: val.CreatedAt,
				UpdatedAt: val.UpdatedAt,
			}
		}

		mockService.getManyMock = func(ctx context.Context, request GetUsersManyRequest) (GetUsersManyResponse, error) {
			assert.EqualValues(t, req, request)

			return expected, nil
		}

		request, err := http.NewRequest(http.MethodGet, "/users", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, request)

		expectedBytes, err := json.Marshal(expected)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, expectedBytes, rr.Body.Bytes())
	})

	t.Run("success with valid query parameters", func(t *testing.T) {
		req := GetUsersManyRequest{
			Page:    1,
			PerPage: 10,
			Filter: User{
				Country: "UK",
			},
		}
		countryParam := fmt.Sprintf("{\"country\":\"%s\"}", req.Filter.Country)

		expected := GetUsersManyResponse{}
		expected.Users = make([]User, len(entities))
		for i, val := range entities {
			expected.Users[i] = User{
				Id:        val.Id,
				FirstName: val.FirstName,
				LastName:  val.LastName,
				Nickname:  val.Nickname,
				Password:  val.Password,
				Email:     val.Email,
				Country:   val.Country,
				CreatedAt: val.CreatedAt,
				UpdatedAt: val.UpdatedAt,
			}
		}

		mockService.getManyMock = func(ctx context.Context, request GetUsersManyRequest) (GetUsersManyResponse, error) {
			assert.EqualValues(t, req, request)

			return expected, nil
		}

		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/users?page=%d&perPage=%d&filter=%v", req.Page, req.PerPage, countryParam), nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, request)

		expectedBytes, err := json.Marshal(expected)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, expectedBytes, rr.Body.Bytes())
	})

	t.Run("should return bad request when page query parameter is in invalid syntax", func(t *testing.T) {
		invalidParam := "s"
		invalidSyntaxErr := &strconv.NumError{Err: strconv.ErrSyntax, Func: "Atoi", Num: invalidParam}
		expected := apierr.BadRequest(invalidSyntaxErr.Error())

		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/users?page=%s", invalidParam), nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, request)

		var actualResp apierr.ApiError
		err = json.Unmarshal(rr.Body.Bytes(), &actualResp)

		assert.Equal(t, expected.StatusCode, rr.Code)
		assert.Equal(t, expected.Code, actualResp.Code)
		assert.Equal(t, expected.Message, actualResp.Message)
	})

	t.Run("should return bad request when perPage query parameter is in invalid syntax", func(t *testing.T) {
		invalidParam := "s"
		invalidSyntaxErr := &strconv.NumError{Err: strconv.ErrSyntax, Func: "Atoi", Num: invalidParam}
		expected := apierr.BadRequest(invalidSyntaxErr.Error())

		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/users?perPage=%s", invalidParam), nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, request)

		var actualResp apierr.ApiError
		err = json.Unmarshal(rr.Body.Bytes(), &actualResp)

		assert.Equal(t, expected.StatusCode, rr.Code)
		assert.Equal(t, expected.Code, actualResp.Code)
		assert.Equal(t, expected.Message, actualResp.Message)
	})

	t.Run("should return bad request when filter parameter schema is invalid", func(t *testing.T) {
		invalidParam := "test"
		expected := apierr.BadRequest("")

		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/users?filter=%s", invalidParam), nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, request)

		var actualResp apierr.ApiError
		err = json.Unmarshal(rr.Body.Bytes(), &actualResp)

		assert.Equal(t, expected.StatusCode, rr.Code)
		assert.Equal(t, expected.Code, actualResp.Code)
	})

	t.Run("should return internal error when service fails", func(t *testing.T) {

		expected := repositoryError(fmt.Errorf("mock error"))
		mockService.getManyMock = func(ctx context.Context, request GetUsersManyRequest) (GetUsersManyResponse, error) {
			return GetUsersManyResponse{}, expected
		}

		request, err := http.NewRequest(http.MethodGet, "/users", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, request)

		var actualResp apierr.ApiError
		err = json.Unmarshal(rr.Body.Bytes(), &actualResp)

		assert.Equal(t, expected.StatusCode, rr.Code)
		assert.Equal(t, expected.Code, actualResp.Code)
		assert.Equal(t, expected.Message, actualResp.Message)
	})
}
