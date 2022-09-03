// Package notify contains operations for notifying external services.
package notify

import (
	"faceit-backend-test/internal/pubsub"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
)

// notificationWebhookHandler subscribes to a certain topic
// after a callback is verified, it is passed to this struct
// it sends the notifications triggered by the topic to the subscribers.
type notificationWebhookHandler struct {
	logger              *logrus.Logger
	eventSubscriber     *pubsub.Subscriber
	notificationClients []*notificationWebhookClient
	mtx                 sync.RWMutex
}

func newTopicHandler(logger *logrus.Logger, subscriber *pubsub.Subscriber) *notificationWebhookHandler {
	return &notificationWebhookHandler{
		logger:              logger,
		eventSubscriber:     subscriber,
		notificationClients: make([]*notificationWebhookClient, 0),
	}
}

func (t *notificationWebhookHandler) addToSubscribers(notificationClient *notificationWebhookClient) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.notificationClients = append(t.notificationClients, notificationClient)
}

// notify consumes the channel of one certain topic
// when something pushed to that topic it sends the notification to the all subscribers
func (t *notificationWebhookHandler) notify() {
	for msg := range t.eventSubscriber.Messages() {
		payload := notificationPayload{
			messageId: uuid.New().String(),
			Data:      msg.Body(),
			CreatedAt: msg.CreatedAt(),
		}

		t.mtx.RLock()
		for _, nc := range t.notificationClients {
			err := nc.sendNotification(payload)
			if err != nil {
				t.logger.WithFields(logrus.Fields{
					"subId":     nc.subscriberDetails.Id,
					"messageId": payload.messageId,
					"error":     err,
				}).Errorln("sending notification to retry queue")
				// TODO: send to retry queue
				continue
			}

			t.logger.WithFields(logrus.Fields{
				"payload":    payload,
				"subscriber": nc.subscriberDetails,
			}).Debug("notification sent successfully")
		}
		t.mtx.RUnlock()
	}
}
