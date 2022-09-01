package notify

import (
	"faceit-backend-test/internal/pubsub"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
)

type topicHandler struct {
	logger                  *logrus.Logger
	eventSubscriber         *pubsub.Subscriber
	notificationSubscribers []NotificationSubscriber
	mtx                     sync.RWMutex
}

func newTopicHandler(logger *logrus.Logger, subscriber *pubsub.Subscriber) *topicHandler {
	return &topicHandler{
		logger:                  logger,
		eventSubscriber:         subscriber,
		notificationSubscribers: make([]NotificationSubscriber, 0),
	}
}

func (t *topicHandler) addToSubscribers(subscriber NotificationSubscriber) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.notificationSubscribers = append(t.notificationSubscribers, subscriber)
}

func (t *topicHandler) notify() {
	for msg := range t.eventSubscriber.Messages() {
		payload := notificationPayload{
			messageId: uuid.New().String(),
			Data:      msg.Body(),
			CreatedAt: msg.CreatedAt(),
		}

		t.mtx.RLock()
		for _, sub := range t.notificationSubscribers {
			err := sendNotification(sub, payload)
			if err != nil {
				t.logger.WithFields(logrus.Fields{
					"subId":     sub.Id,
					"messageId": payload.messageId,
					"error":     err,
				}).Errorln("sending notification to retry queue")
				// TODO: send to retry queue
				continue
			}

			t.logger.WithFields(logrus.Fields{
				"payload":    payload,
				"subscriber": sub,
			}).Debug("notification sent successfully")
		}
		t.mtx.RUnlock()
	}
}
