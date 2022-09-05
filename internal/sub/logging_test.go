package sub

import (
	"context"
	"faceit-backend-test/internal/apierr"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestServiceLoggingMiddleware_Subscribe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := SubscribeRequest{
			Type:     "test",
			Callback: "http://callback.com",
			Secret:   "xxx",
		}

		expected := SubscribeResponse{
			Id:        uuid.New().String(),
			Type:      "test",
			Status:    "",
			CreatedAt: time.Now(),
		}

		mockService := &mockService{}
		mockService.subscribeMock = func(ctx context.Context, request SubscribeRequest) (SubscribeResponse, error) {
			assert.EqualValues(t, req, request)
			return expected, nil
		}

		logger := logrus.New()
		loggingMiddleware := NewServiceLoggingMiddleware(logger)(mockService)

		resp, err := loggingMiddleware.Subscribe(context.Background(), req)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, resp)
	})

	t.Run("error", func(t *testing.T) {
		req := SubscribeRequest{
			Type:     "test",
			Callback: "http://callback.com",
			Secret:   "xxx",
		}

		expectedErr := SubscribeError("mock error")

		mockService := &mockService{}
		mockService.subscribeMock = func(ctx context.Context, request SubscribeRequest) (SubscribeResponse, error) {
			assert.EqualValues(t, req, request)
			return SubscribeResponse{}, expectedErr
		}

		logger := logrus.New()
		loggingMiddleware := NewServiceLoggingMiddleware(logger)(mockService)

		_, err := loggingMiddleware.Subscribe(context.Background(), req)
		assert.Error(t, err)
		assert.IsType(t, expectedErr, err)

		apiErr := err.(apierr.ApiError)
		assert.EqualValues(t, expectedErr, apiErr)
	})
}
