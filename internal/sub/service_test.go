package sub

import (
	"context"
	"faceit-backend-test/internal/apierr"
	"faceit-backend-test/internal/notify"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockNotificationManager struct {
	subscribeMock func(topic string, callback string, secret string) (notify.NotificationSubscriber, error)
}

func (m *mockNotificationManager) Subscribe(topic string, callback string, secret string) (notify.NotificationSubscriber, error) {
	return m.subscribeMock(topic, callback, secret)
}

func TestService_Subscribe(t *testing.T) {
	mockNotifMngr := &mockNotificationManager{}
	service := NewService(WithNotificationManager(mockNotifMngr))

	t.Run("success", func(t *testing.T) {
		mockResponse := notify.NotificationSubscriber{
			Id:          uuid.New().String(),
			Topic:       "test",
			CallbackUrl: "http://url.com",
			Secret:      "verysecret",
			Status:      "status",
			CreatedAt:   time.Now(),
		}

		expected := SubscribeResponse{
			Id:        mockResponse.Id,
			Type:      mockResponse.Topic,
			Status:    mockResponse.Status,
			CreatedAt: mockResponse.CreatedAt,
		}
		req := SubscribeRequest{
			Type:     mockResponse.Topic,
			Callback: mockResponse.CallbackUrl,
			Secret:   mockResponse.Secret,
		}

		mockNotifMngr.subscribeMock = func(topic string, callback string, secret string) (notify.NotificationSubscriber, error) {
			assert.Equal(t, req.Type, topic)
			assert.Equal(t, req.Callback, callback)
			assert.Equal(t, req.Secret, secret)

			return mockResponse, nil
		}

		actual, err := service.Subscribe(context.Background(), req)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, actual)
	})

	t.Run("error", func(t *testing.T) {
		expected := fmt.Errorf("mock error")
		expectedApiErr := SubscribeError(expected.Error())
		mockNotifMngr.subscribeMock = func(topic string, callback string, secret string) (notify.NotificationSubscriber, error) {
			return notify.NotificationSubscriber{}, expected
		}

		req := SubscribeRequest{}
		_, err := service.Subscribe(context.Background(), req)
		assert.Error(t, err)
		assert.IsType(t, expectedApiErr, err)

		apiErr := err.(apierr.ApiError)
		assert.EqualValues(t, expectedApiErr, apiErr)
	})

}
