package notify

import (
	"faceit-backend-test/internal/pubsub"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type NotificationManager interface {
	Subscribe(topic string, callback string, secret string) (NotificationSubscriber, error)
}

type notificationManager struct {
	broker              *pubsub.Broker
	pendingVerification chan NotificationSubscriber
	topicHandlers       map[string]*topicHandler
	logger              *logrus.Logger
}

var _ NotificationManager = (*notificationManager)(nil)

type NotificationManagerOpts func(*notificationManager)

func NewNotificationManager(topics []string, opts ...NotificationManagerOpts) *notificationManager {
	n := &notificationManager{
		topicHandlers:       map[string]*topicHandler{},
		pendingVerification: make(chan NotificationSubscriber, 100),
	}

	for _, opt := range opts {
		opt(n)
	}

	for _, topic := range topics {
		subscriber := pubsub.NewSubscriber()
		n.broker.Subscribe(subscriber, topic)
		n.topicHandlers[topic] = newTopicHandler(n.logger, subscriber)
	}

	return n
}

func WithLogger(logger *logrus.Logger) NotificationManagerOpts {
	return func(manager *notificationManager) {
		manager.logger = logger
	}
}

func WithBroker(broker *pubsub.Broker) NotificationManagerOpts {
	return func(manager *notificationManager) {
		manager.broker = broker
	}
}

func (n *notificationManager) Start() {
	go n.handlePendingVerifications()

	for _, handler := range n.topicHandlers {
		go handler.notify()
	}
}

func (n *notificationManager) Subscribe(topic string, callback string, secret string) (NotificationSubscriber, error) {
	_, ok := n.topicHandlers[topic]
	if !ok {
		return NotificationSubscriber{}, fmt.Errorf("topic %s not available", topic)
	}

	ns := NotificationSubscriber{
		Id:          uuid.New().String(),
		Topic:       topic,
		CallbackUrl: callback,
		Secret:      secret,
		Status:      verificationStatusPending,
		CreatedAt:   time.Now(),
	}

	n.logger.WithFields(logrus.Fields{
		"id":       ns.Id,
		"topic":    ns.Topic,
		"callback": ns.CallbackUrl,
		"status":   ns.Status,
	}).Debug("pushed to pending verifications")
	n.pendingVerification <- ns
	return ns, nil
}

func (n *notificationManager) handlePendingVerifications() {
	for sub := range n.pendingVerification {
		err := verifySubscription(sub)
		if err != nil {
			n.logger.WithFields(logrus.Fields{
				"id":       sub.Id,
				"topic":    sub.Topic,
				"callback": sub.CallbackUrl,
				"error":    err,
			}).Error("cannot verify subscription")
			continue
		}

		topic, ok := n.topicHandlers[sub.Topic]
		if !ok {
			continue
		}

		n.logger.WithFields(logrus.Fields{
			"id":       sub.Id,
			"topic":    sub.Topic,
			"callback": sub.CallbackUrl,
		}).Debug("verified subscription")
		topic.addToSubscribers(sub)
	}
}
