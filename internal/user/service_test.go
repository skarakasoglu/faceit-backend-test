package user

import (
	"context"
	"faceit-backend-test/internal/apierr"
	"faceit-backend-test/internal/pubsub"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type mockRepository struct {
	createMock     func(context.Context, Entity) (Entity, error)
	updateMock     func(context.Context, Entity) (Entity, error)
	deleteByIdMock func(context.Context, string) error
	getManyMock    func(context.Context, GetManyParameters) ([]Entity, error)
}

func (m *mockRepository) Create(ctx context.Context, entity Entity) (Entity, error) {
	return m.createMock(ctx, entity)
}

func (m *mockRepository) Update(ctx context.Context, entity Entity) (Entity, error) {
	return m.updateMock(ctx, entity)
}

func (m *mockRepository) DeleteById(ctx context.Context, id string) error {
	return m.deleteByIdMock(ctx, id)
}

func (m *mockRepository) GetMany(ctx context.Context, parameters GetManyParameters) ([]Entity, error) {
	return m.getManyMock(ctx, parameters)
}

func TestService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := CreateUserRequest{
			FirstName: e.FirstName,
			LastName:  e.LastName,
			Nickname:  e.Nickname,
			Password:  e.Password,
			Email:     e.Email,
			Country:   e.Country,
		}

		mockRepo := &mockRepository{}
		mockRepo.createMock = func(ctx context.Context, entity Entity) (Entity, error) {
			assert.Equal(t, req.FirstName, entity.FirstName)
			assert.Equal(t, req.LastName, entity.LastName)
			assert.Equal(t, req.Nickname, entity.Nickname)
			assert.Equal(t, req.Password, entity.Password)
			assert.Equal(t, req.Email, entity.Email)
			assert.Equal(t, req.Country, entity.Country)

			return e, nil
		}

		service := NewService(WithRepository(mockRepo))

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

		actual, err := service.Create(context.Background(), req)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, actual)
	})

	t.Run("repository error", func(t *testing.T) {
		expectedErr := fmt.Errorf("mock error")

		mockRepo := &mockRepository{}
		mockRepo.createMock = func(ctx context.Context, entity Entity) (Entity, error) {
			return Entity{}, expectedErr
		}

		service := NewService(WithRepository(mockRepo))

		req := CreateUserRequest{}
		_, err := service.Create(context.Background(), req)

		assert.Error(t, err)
		assert.IsType(t, apierr.ApiError{}, err)

		apiErr := err.(apierr.ApiError)
		assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
		assert.Equal(t, expectedErr.Error(), apiErr.Message)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := UpdateUserRequest{
			Id:        e.Id,
			FirstName: e.FirstName,
			LastName:  e.LastName,
			Nickname:  e.Nickname,
			Password:  e.Password,
			Email:     e.Email,
			Country:   e.Country,
		}

		mockRepo := &mockRepository{}
		mockRepo.updateMock = func(ctx context.Context, entity Entity) (Entity, error) {
			assert.Equal(t, req.Id, entity.Id)
			assert.Equal(t, req.FirstName, entity.FirstName)
			assert.Equal(t, req.LastName, entity.LastName)
			assert.Equal(t, req.Nickname, entity.Nickname)
			assert.Equal(t, req.Password, entity.Password)
			assert.Equal(t, req.Email, entity.Email)
			assert.Equal(t, req.Country, entity.Country)

			return e, nil
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

		mockBroker := &pubsub.MockBroker{}
		mockBroker.PublishMock = func(s string, i interface{}) {
			assert.Equal(t, UserChangeTopic, s)
			assert.IsType(t, User{}, i)

			val := i.(User)
			assert.EqualValues(t, expected.User, val)
		}

		service := NewService(WithRepository(mockRepo), WithBroker(mockBroker))
		actual, err := service.Update(context.Background(), req)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, actual)
	})

	t.Run("repository error", func(t *testing.T) {
		expectedErr := fmt.Errorf("mock error")

		mockRepo := &mockRepository{}
		mockRepo.updateMock = func(ctx context.Context, entity Entity) (Entity, error) {
			return Entity{}, expectedErr
		}

		service := NewService(WithRepository(mockRepo))

		req := UpdateUserRequest{}
		_, err := service.Update(context.Background(), req)

		assert.Error(t, err)
		assert.IsType(t, apierr.ApiError{}, err)

		apiErr := err.(apierr.ApiError)
		assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
		assert.Equal(t, expectedErr.Error(), apiErr.Message)
	})
}

func TestService_DeleteById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := DeleteUserByIdRequest{Id: e.Id}
		expected := DeleteUserResponse{Id: e.Id}

		mockRepo := &mockRepository{}
		mockRepo.deleteByIdMock = func(ctx context.Context, s string) error {
			assert.Equal(t, req.Id, s)

			return nil
		}

		service := NewService(WithRepository(mockRepo))

		actual, err := service.DeleteById(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expected.Id, actual.Id)
	})

	t.Run("repository error", func(t *testing.T) {
		expectedErr := fmt.Errorf("mock error")

		mockRepo := &mockRepository{}
		mockRepo.deleteByIdMock = func(ctx context.Context, id string) error {
			return expectedErr
		}

		service := NewService(WithRepository(mockRepo))

		req := DeleteUserByIdRequest{}
		_, err := service.DeleteById(context.Background(), req)

		assert.Error(t, err)
		assert.IsType(t, apierr.ApiError{}, err)

		apiErr := err.(apierr.ApiError)
		assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
		assert.Equal(t, expectedErr.Error(), apiErr.Message)
	})
}

func TestService_GetMany(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := GetUsersManyRequest{
			Page:    1,
			PerPage: 10,
			Filter:  User{},
		}

		mockRepo := &mockRepository{}
		mockRepo.getManyMock = func(ctx context.Context, parameters GetManyParameters) ([]Entity, error) {
			assert.Equal(t, req.Page, parameters.Page)
			assert.Equal(t, req.PerPage, parameters.PerPage)
			assert.Equal(t, req.Filter.Id, parameters.Filter.Id)
			assert.Equal(t, req.Filter.FirstName, parameters.Filter.FirstName)
			assert.Equal(t, req.Filter.LastName, parameters.Filter.LastName)
			assert.Equal(t, req.Filter.Nickname, parameters.Filter.Nickname)
			assert.Equal(t, req.Filter.Password, parameters.Filter.Password)
			assert.Equal(t, req.Filter.Country, parameters.Filter.Country)
			assert.Equal(t, req.Filter.CreatedAt, parameters.Filter.CreatedAt)
			assert.Equal(t, req.Filter.UpdatedAt, parameters.Filter.UpdatedAt)

			return entities, nil
		}

		expected := GetUsersManyResponse{}
		expected.Users = make([]User, len(entities))
		for i, e := range entities {
			expected.Users[i] = User{
				Id:        e.Id,
				FirstName: e.FirstName,
				LastName:  e.LastName,
				Nickname:  e.Nickname,
				Password:  e.Password,
				Email:     e.Email,
				Country:   e.Country,
				CreatedAt: e.CreatedAt,
				UpdatedAt: e.UpdatedAt,
			}
		}

		service := NewService(WithRepository(mockRepo))
		actual, err := service.GetMany(context.Background(), req)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, actual)
	})

	t.Run("repository error", func(t *testing.T) {
		expectedErr := fmt.Errorf("mock error")

		mockRepo := &mockRepository{}
		mockRepo.getManyMock = func(ctx context.Context, parameters GetManyParameters) ([]Entity, error) {
			return []Entity{}, expectedErr
		}

		service := NewService(WithRepository(mockRepo))

		req := GetUsersManyRequest{}
		_, err := service.GetMany(context.Background(), req)

		assert.Error(t, err)
		assert.IsType(t, apierr.ApiError{}, err)

		apiErr := err.(apierr.ApiError)
		assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
		assert.Equal(t, expectedErr.Error(), apiErr.Message)
	})
}
