package notify

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// newTestClient returns *http.Client with Transport replaced to avoid making real calls
func newTestClient(fn RoundTripFunc, subscriber NotificationSubscriber) *notificationWebhookClient {
	return &notificationWebhookClient{
		client: &http.Client{
			Transport: fn,
		},
		subscriberDetails: subscriber,
	}
}

func TestNotificationWebhookClient_Verify(t *testing.T) {
	ns := NotificationSubscriber{
		Id:          uuid.New().String(),
		Topic:       "test",
		CallbackUrl: "http://example.com",
		Secret:      "",
		CreatedAt:   time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		nc := newTestClient(func(req *http.Request) *http.Response {
			messageType := req.Header.Get("Notification-Message-Type")
			assert.Equal(t, notificationMessageTypeVerification, messageType)

			buffer, err := io.ReadAll(req.Body)
			assert.NoError(t, err)

			var payload verificationPayload
			err = json.Unmarshal(buffer, &payload)
			assert.NoError(t, err)

			assert.Equal(t, ns.CallbackUrl, payload.Callback)
			assert.Equal(t, ns.Id, payload.Id)
			assert.Equal(t, ns.CreatedAt.Unix(), payload.CreatedAt.Unix())

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(payload.Challenge)),
			}
		}, ns)

		err := nc.Verify()
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		nc := newTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
			}
		}, ns)

		err := nc.Verify()
		assert.Error(t, err)
	})
}

func TestNotificationWebhookClient_Notify(t *testing.T) {

	ns := NotificationSubscriber{
		Id:          uuid.New().String(),
		Topic:       "test",
		CallbackUrl: "http://example.com",
		Secret:      "",
		CreatedAt:   time.Now(),
	}

	t.Run("success with no secret provided", func(t *testing.T) {
		notification := notificationPayload{
			messageId: uuid.New().String(),
			timestamp: time.Now(),
			Data: struct {
				Test int `json:"test"`
			}{5},
			CreatedAt: time.Now(),
		}

		nc := newTestClient(func(req *http.Request) *http.Response {
			messageType := req.Header.Get("Notification-Message-Type")
			assert.Equal(t, notificationMessageTypeNotification, messageType)
			assert.Equal(t, notification.messageId, req.Header.Get("Notification-Message-Id"))
			assert.NotEqual(t, "", req.Header.Get("Notification-Message-Timestamp"))

			buffer, err := io.ReadAll(req.Body)
			assert.NoError(t, err)

			expectedBuffer, err := json.Marshal(notification)
			assert.NoError(t, err)

			assert.Equal(t, expectedBuffer, buffer)

			return &http.Response{
				StatusCode: http.StatusOK,
			}
		}, ns)

		err := nc.Notify(notification)
		assert.NoError(t, err)
	})

	t.Run("success with secret provided", func(t *testing.T) {
		ns.Secret = "ABCdE56"

		notification := notificationPayload{
			messageId: uuid.New().String(),
			timestamp: time.Now(),
			Data: struct {
				Test int `json:"test"`
			}{5},
			CreatedAt: time.Now(),
		}

		nc := newTestClient(func(req *http.Request) *http.Response {
			messageType := req.Header.Get("Notification-Message-Type")
			assert.Equal(t, notificationMessageTypeNotification, messageType)

			messageId := req.Header.Get("Notification-Message-Id")
			assert.Equal(t, notification.messageId, messageId)

			messageTimestampStr := req.Header.Get("Notification-Message-Timestamp")
			assert.NotEqual(t, "", messageTimestampStr)

			buffer, err := io.ReadAll(req.Body)
			assert.NoError(t, err)

			messageTimestamp, err := time.Parse(time.RFC3339Nano, messageTimestampStr)
			assert.NoError(t, err)

			h := hmac.New(sha256.New, []byte(ns.Secret))
			h.Write([]byte(fmt.Sprintf("%v%v%v", messageId, messageTimestamp.Unix(), string(buffer))))
			expectedSignature := fmt.Sprintf("sha256=%v", hex.EncodeToString(h.Sum(nil)))

			signature := req.Header.Get("Notification-Message-Signature")
			assert.Equal(t, expectedSignature, signature)

			expectedBuffer, err := json.Marshal(notification)
			assert.NoError(t, err)

			assert.Equal(t, expectedBuffer, buffer)

			return &http.Response{
				StatusCode: http.StatusOK,
			}
		}, ns)

		err := nc.Notify(notification)
		assert.NoError(t, err)
	})
}
