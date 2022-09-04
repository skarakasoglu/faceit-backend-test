package sub

import (
	"bytes"
	"context"
	"encoding/json"
	"faceit-backend-test/internal/apierr"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockService struct {
	subscribeMock func(ctx context.Context, request SubscribeRequest) (SubscribeResponse, error)
}

func (m *mockService) Subscribe(ctx context.Context, request SubscribeRequest) (SubscribeResponse, error) {
	return m.subscribeMock(ctx, request)
}

func TestController_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	controller := NewController()

	t.Run("", func(t *testing.T) {
		controller.Register(&router.RouterGroup)

		rr := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, "/subscribe", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.NotEqual(t, http.StatusNotFound, rr.Code)
	})
}

func TestController_Subscribe(t *testing.T) {
	mockService := &mockService{}
	gin.SetMode(gin.TestMode)
	controller := NewController(WithService(mockService))
	router := gin.Default()
	router.POST(route, controller.Subscribe)

	t.Run("success", func(t *testing.T) {
		expected := SubscribeResponse{
			Id:        uuid.New().String(),
			Type:      "test",
			Status:    "status",
			CreatedAt: time.Now(),
		}
		req := SubscribeRequest{
			Type:     expected.Type,
			Callback: "http://callback.com",
		}

		mockService.subscribeMock = func(ctx context.Context, request SubscribeRequest) (SubscribeResponse, error) {
			assert.EqualValues(t, req, request)

			return expected, nil
		}

		reqBody, err := json.Marshal(req)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(expected)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
	})

	t.Run("should return bad request when a required field is missing", func(t *testing.T) {

		req := SubscribeRequest{
			Type:   "test",
			Secret: "xadsad",
		}
		expected := apierr.BadRequest("invalid request body")

		reqBody, err := json.Marshal(req)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, request)

		var actualResp apierr.ApiError
		err = json.Unmarshal(rr.Body.Bytes(), &actualResp)

		assert.Equal(t, expected.StatusCode, rr.Code)
		assert.Equal(t, expected.Code, actualResp.Code)
	})

	t.Run("should return internal error when service fails", func(t *testing.T) {
		req := SubscribeRequest{
			Type:     "test",
			Callback: "callback",
			Secret:   "xadsad",
		}

		expected := SubscribeError("mock error")
		mockService.subscribeMock = func(ctx context.Context, request SubscribeRequest) (SubscribeResponse, error) {
			return SubscribeResponse{}, expected
		}

		reqBody, err := json.Marshal(req)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(reqBody))
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
