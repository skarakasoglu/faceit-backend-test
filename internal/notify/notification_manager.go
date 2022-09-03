package notify

import (
	"faceit-backend-test/internal/pubsub"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type NotificationManager interface {
	Subscribe(topic string, callback string, secret string) (NotificationSubscriber, error)
}

// notificationManager verifies the subscriptions and
// creates the notification handlers topic by topic
type notificationManager struct {
	broker                      *pubsub.Broker
	pendingVerification         chan *notificationWebhookClient
	notificationWebhookHandlers map[string]*notificationWebhookHandler
	logger                      *logrus.Logger
}

var _ NotificationManager = (*notificationManager)(nil)

type NotificationManagerOpts func(*notificationManager)

func NewNotificationManager(topics []string, opts ...NotificationManagerOpts) *notificationManager {
	n := &notificationManager{
		notificationWebhookHandlers: map[string]*notificationWebhookHandler{},
		pendingVerification:         make(chan *notificationWebhookClient, 100),
	}

	for _, opt := range opts {
		opt(n)
	}

	for _, topic := range topics {
		subscriber := pubsub.NewSubscriber()
		n.broker.Subscribe(subscriber, topic)
		n.notificationWebhookHandlers[topic] = newTopicHandler(n.logger, subscriber)
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

// Start starts all the notification operators async
func (n *notificationManager) Start() {
	go n.handlePendingVerifications()

	for _, handler := range n.notificationWebhookHandlers {
		go handler.notify()
	}
}

// Subscribe creates a new subscriber and sends it to the verification queue
func (n *notificationManager) Subscribe(topic string, callback string, secret string) (NotificationSubscriber, error) {
	_, ok := n.notificationWebhookHandlers[topic]
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
	nwc := newNotificationWebhookClient(&http.Client{}, ns)

	n.logger.WithFields(logrus.Fields{
		"id":       ns.Id,
		"topic":    ns.Topic,
		"callback": ns.CallbackUrl,
		"status":   ns.Status,
	}).Debug("pushed to pending verifications")
	n.pendingVerification <- nwc
	return ns, nil
}

// handlePendingVerifications checks pendingVerification channel regularly
// if something is pushed then sends a verification request to the callback url of the subscriber.
func (n *notificationManager) handlePendingVerifications() {
	for notifierClient := range n.pendingVerification {
		err := notifierClient.verifySubscription()
		if err != nil {
			n.logger.WithFields(logrus.Fields{
				"id":       notifierClient.subscriberDetails.Id,
				"topic":    notifierClient.subscriberDetails.Topic,
				"callback": notifierClient.subscriberDetails.CallbackUrl,
				"error":    err,
			}).Error("cannot verify subscription")
			continue
		}

		topic, ok := n.notificationWebhookHandlers[notifierClient.subscriberDetails.Topic]
		if !ok {
			continue
		}

		n.logger.WithFields(logrus.Fields{
			"id":       notifierClient.subscriberDetails.Id,
			"topic":    notifierClient.subscriberDetails.Topic,
			"callback": notifierClient.subscriberDetails.CallbackUrl,
		}).Debug("verified subscription")
		topic.addToSubscribers(notifierClient)
	}
}
