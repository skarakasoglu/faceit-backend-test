package notify

import "time"

const (
	verificationStatusPending = "verification_pending"
	verificationStatusEnabled = "enabled"

	notificationMessageTypeVerification = "webhook_callback_verification"
	notificationMessageTypeNotification = "notification"
)

type NotificationSubscriber struct {
	Id          string
	Topic       string
	CallbackUrl string
	Secret      string
	Status      string
	CreatedAt   time.Time
}

type verificationPayload struct {
	Id        string    `json:"id"`
	Challenge string    `json:"challenge"`
	Callback  string    `json:"callback"`
	CreatedAt time.Time `json:"createdAt"`
}

type notificationPayload struct {
	messageId string
	timestamp time.Time
	Data      interface{} `json:"data"`
	CreatedAt time.Time   `json:"createdAt"`
}
