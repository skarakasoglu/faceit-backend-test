package notify

import (
	"faceit-backend-test/internal/pubsub"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
