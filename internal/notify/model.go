package notify

import "time"

const (
	verificationStatusPending = "verification_pending"
	verificationStatusEnabled = "enabled"

	notificationMessageTypeVerification = "webhook_callback_verification"
	notificationMessageTypeNotification = "notification"
)

// NotificationSubscriber is the model used for storing subscriber information
// this information is passed to the subscriber in the initial request.
type NotificationSubscriber struct {
	Id          string
	Topic       string
	CallbackUrl string
	Secret      string
	Status      string
	CreatedAt   time.Time
}

// verificationPayload is the model used in verification flow of the webhook.
type verificationPayload struct {
	Id        string    `json:"id"`
	Challenge string    `json:"challenge"`
	Callback  string    `json:"callback"`
	CreatedAt time.Time `json:"createdAt"`
}

// notificationPayload is the model used in notification flow of the webhook.
type notificationPayload struct {
	messageId string
	timestamp time.Time
	Data      interface{} `json:"data"`
	CreatedAt time.Time   `json:"createdAt"`
}
