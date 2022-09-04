package user

import (
	"context"
	"faceit-backend-test/internal/apierr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestServiceLoggingMiddleware_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		serviceMock := &mockService{}

		req := CreateUserRequest{
			FirstName: e.FirstName,
			LastName:  e.LastName,
			Nickname:  e.Nickname,
			Password:  e.Password,
			Email:     e.Email,
			Country:   e.Country,
		}
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
		serviceMock.createMock = func(ctx context.Context, request CreateUserRequest) (CreateUserResponse, error) {
			assert.EqualValues(t, req, request)

			return expected, nil
		}

		logger := logrus.New()
		loggingMiddleware := NewServiceLoggingMiddleware(logger)(serviceMock)

		resp, err := loggingMiddleware.Create(context.Background(), req)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, resp)
	})

	t.Run("error", func(t *testing.T) {
		serviceMock := &mockService{}
		expected := apierr.BadRequest("")
		req := CreateUserRequest{}

		serviceMock.createMock = func(ctx context.Context, request CreateUserRequest) (CreateUserResponse, error) {
			return CreateUserResponse{}, expected
		}

		logger := logrus.New()
		loggingMiddleware := NewServiceLoggingMiddleware(logger)(serviceMock)

		_, err := loggingMiddleware.Create(context.Background(), req)
		assert.Error(t, err)
		assert.IsType(t, expected, err)

		apiErr := err.(apierr.ApiError)
		assert.EqualValues(t, expected, apiErr)
	})
}

func TestServiceLoggingMiddleware_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		serviceMock := &mockService{}

		req := UpdateUserRequest{
			Id:        e.Id,
			FirstName: e.FirstName,
			LastName:  e.LastName,
			Nickname:  e.Nickname,
			Password:  e.Password,
			Email:     e.Email,
			Country:   e.Country,
		}
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
		serviceMock.updateMock = func(ctx context.Context, request UpdateUserRequest) (UpdateUserResponse, error) {
			assert.EqualValues(t, req, request)

			return expected, nil
		}

		logger := logrus.New()
		loggingMiddleware := NewServiceLoggingMiddleware(logger)(serviceMock)

		resp, err := loggingMiddleware.Update(context.Background(), req)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, resp)
	})

	t.Run("error", func(t *testing.T) {
		serviceMock := &mockService{}
		expected := apierr.BadRequest("")
		req := UpdateUserRequest{}

		serviceMock.updateMock = func(ctx context.Context, request UpdateUserRequest) (UpdateUserResponse, error) {
			return UpdateUserResponse{}, expected
		}

		logger := logrus.New()
		loggingMiddleware := NewServiceLoggingMiddleware(logger)(serviceMock)

		_, err := loggingMiddleware.Update(context.Background(), req)
		assert.Error(t, err)
		assert.IsType(t, expected, err)

		apiErr := err.(apierr.ApiError)
		assert.EqualValues(t, expected, apiErr)
	})
}

func TestServiceLoggingMiddleware_DeleteById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		serviceMock := &mockService{}

		req := DeleteUserByIdRequest{
			Id: e.Id,
		}
		expected := DeleteUserResponse{
			Id: e.Id,
		}
		serviceMock.deleteByIdMock = func(ctx context.Context, request DeleteUserByIdRequest) (DeleteUserResponse, error) {
			assert.EqualValues(t, req, request)

			return expected, nil
		}

		logger := logrus.New()
		loggingMiddleware := NewServiceLoggingMiddleware(logger)(serviceMock)

		resp, err := loggingMiddleware.DeleteById(context.Background(), req)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, resp)
	})

	t.Run("error", func(t *testing.T) {
		serviceMock := &mockService{}
		expected := apierr.BadRequest("")
		req := DeleteUserByIdRequest{}

		serviceMock.deleteByIdMock = func(ctx context.Context, request DeleteUserByIdRequest) (DeleteUserResponse, error) {
			return DeleteUserResponse{}, expected
		}

		logger := logrus.New()
		loggingMiddleware := NewServiceLoggingMiddleware(logger)(serviceMock)

		_, err := loggingMiddleware.DeleteById(context.Background(), req)
		assert.Error(t, err)
		assert.IsType(t, expected, err)

		apiErr := err.(apierr.ApiError)
		assert.EqualValues(t, expected, apiErr)
	})
}

func TestServiceLoggingMiddleware_GetMany(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		serviceMock := &mockService{}

		req := GetUsersManyRequest{
			Page:    3,
			PerPage: 25,
			Filter: User{
				Id:        "",
				FirstName: "",
				LastName:  "",
				Nickname:  "",
				Password:  "",
				Email:     "",
				Country:   "UK",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
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

		serviceMock.getManyMock = func(ctx context.Context, request GetUsersManyRequest) (GetUsersManyResponse, error) {
			assert.EqualValues(t, req, request)

			return expected, nil
		}

		logger := logrus.New()
		loggingMiddleware := NewServiceLoggingMiddleware(logger)(serviceMock)

		resp, err := loggingMiddleware.GetMany(context.Background(), req)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, resp)
	})

	t.Run("error", func(t *testing.T) {
		serviceMock := &mockService{}
		expected := apierr.BadRequest("")
		req := GetUsersManyRequest{}

		serviceMock.getManyMock = func(ctx context.Context, request GetUsersManyRequest) (GetUsersManyResponse, error) {
			return GetUsersManyResponse{}, expected
		}

		logger := logrus.New()
		loggingMiddleware := NewServiceLoggingMiddleware(logger)(serviceMock)

		_, err := loggingMiddleware.GetMany(context.Background(), req)
		assert.Error(t, err)
		assert.IsType(t, expected, err)

		apiErr := err.(apierr.ApiError)
		assert.EqualValues(t, expected, apiErr)
	})
}
