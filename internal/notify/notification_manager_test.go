package notify

import (
	"faceit-backend-test/internal/pubsub"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

type mockNotificationWebhookClient struct {
	subscriber NotificationSubscriber
	verifyMock func() error
	notifyMock func(notificationPayload) error
}

func (m *mockNotificationWebhookClient) Verify() error {
	return m.verifyMock()
}

func (m *mockNotificationWebhookClient) Notify(payload notificationPayload) error {
	return m.notifyMock(payload)
}

func (m *mockNotificationWebhookClient) Subscriber() NotificationSubscriber {
	return m.subscriber
}

func TestNotificationManager_Subscribe(t *testing.T) {
	logger := logrus.New()
	broker := &pubsub.MockBroker{}
	topics := []string{"test"}
	broker.SubscribeMock = func(subscriber pubsub.Subscriber, s string) {
	}
	ntfMngr := NewNotificationManager(topics, WithLogger(logger), WithBroker(broker))

	t.Run("success", func(t *testing.T) {

		topic := "test"
		callback := "callback"
		secret := "supersecret"

		resp, err := ntfMngr.Subscribe(topic, callback, secret)
		assert.NoError(t, err)
		assert.Equal(t, topic, resp.Topic)
		assert.Equal(t, callback, resp.CallbackUrl)
		assert.Equal(t, secret, resp.Secret)
	})

	t.Run("should return error when topic doesn't exist", func(t *testing.T) {
		nonExistingTopic := "xxxx"

		_, err := ntfMngr.Subscribe(nonExistingTopic, "", "")
		assert.Error(t, err)
	})
}

func TestNotificationManager_handlePendingVerifications(t *testing.T) {
	logger := logrus.New()
	broker := &pubsub.MockBroker{}
	topics := []string{"test"}
	broker.SubscribeMock = func(subscriber pubsub.Subscriber, s string) {
	}
	ntfMngr := NewNotificationManager(topics, WithLogger(logger), WithBroker(broker))

	t.Run("success", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			ntfMngr.handlePendingVerifications()
			wg.Done()
		}()

		ns := NotificationSubscriber{
			Id:          uuid.New().String(),
			Topic:       "test",
			CallbackUrl: "http://example.com",
			Secret:      "",
			Status:      verificationStatusPending,
			CreatedAt:   time.Now(),
		}
		mockClient := &mockNotificationWebhookClient{subscriber: ns}
		mockClient.verifyMock = func() error {
			return nil
		}

		ntfMngr.pendingVerification <- mockClient
		close(ntfMngr.pendingVerification)
		wg.Wait()
	})

	t.Run("error", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			ntfMngr.handlePendingVerifications()
			wg.Done()
		}()

		ns := NotificationSubscriber{
			Id:          uuid.New().String(),
			Topic:       "test",
			CallbackUrl: "http://example.com",
			Secret:      "",
			Status:      verificationStatusPending,
			CreatedAt:   time.Now(),
		}
		mockClient := &mockNotificationWebhookClient{subscriber: ns}
		mockClient.verifyMock = func() error {
			return fmt.Errorf("mock error")
		}

		ntfMngr.pendingVerification = make(chan NotificationWebhookClient, 10)
		ntfMngr.pendingVerification <- mockClient
		close(ntfMngr.pendingVerification)
		wg.Wait()
	})
}
