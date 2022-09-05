package notify

import (
	"faceit-backend-test/internal/pubsub"
	"github.com/sirupsen/logrus"
	"sync"
	"testing"
)

func TestNotificationWebhookHandler_notify(t *testing.T) {

	logger := logrus.New()
	subscriber := pubsub.NewSubscriber()

	notifyHandler := newNotificationWebhookHandler(logger, subscriber)
	t.Run("success", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			notifyHandler.notify()
			wg.Done()
		}()

		mockClient := &mockNotificationWebhookClient{}
		mockClient.notifyMock = func(payload notificationPayload) error {
			return nil
		}
		notifyHandler.addToSubscribers(mockClient)

		subscriber.Signal(pubsub.NewMessage(5, "test"))
		subscriber.Destruct()
		wg.Wait()
	})
}
