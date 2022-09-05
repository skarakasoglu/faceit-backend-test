package notify

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

// notificationWebhookClient handles the subscription verification and sending the notifications to the subscriber.
type notificationWebhookClient struct {
	client            *http.Client
	subscriberDetails NotificationSubscriber
}

func newNotificationWebhookClient(client *http.Client, sub NotificationSubscriber) *notificationWebhookClient {
	return &notificationWebhookClient{
		client:            client,
		subscriberDetails: sub,
	}
}

// verifySubscription the first flow of the subscription is verifying that the callback endpoint is valid.
// this func makes a request to the callback url with a challenge value
// the client should verify it by returning a http.StatusOK response
// with the challenge as a plain text in the body.
func (n *notificationWebhookClient) verifySubscription() error {
	payload := verificationPayload{
		Id:        n.subscriberDetails.Id,
		Challenge: uuid.New().String(),
		Callback:  n.subscriberDetails.CallbackUrl,
		CreatedAt: n.subscriberDetails.CreatedAt,
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	reqBody := bytes.NewReader(buf)
	req, err := http.NewRequest(http.MethodPost, payload.Callback, reqBody)
	if err != nil {
		return err
	}

	// subscriber can distinguish whether the request type is
	// verification or notification by Notification-Message-Type header
	req.Header.Set("Notification-Message-Type", notificationMessageTypeVerification)
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		challenge := string(body)
		// if the challenge provider by the subscriber is not equal to the sent one
		// then the subscription won't be verified.
		if challenge != payload.Challenge {
			return fmt.Errorf("invalid challenge provided: %v, expected: %v", challenge, payload.Challenge)
		}
	}

	// after verifying that the callback url is valid
	// then set the subscription status as enabled
	n.subscriberDetails.Status = verificationStatusEnabled
	return nil
}

// sendNotification sends the notifications to the callback url which is provided by the subscriber
func (n *notificationWebhookClient) sendNotification(payload notificationPayload) error {
	buf, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	reqBody := bytes.NewReader(buf)
	req, err := http.NewRequest(http.MethodPost, n.subscriberDetails.CallbackUrl, reqBody)
	if err != nil {
		return err
	}

	payload.timestamp = time.Now()

	// message id can be used by the subscriber to be able to
	// handle a possible duplicate notifications from Notification-Message-Id header in the request.
	req.Header.Set("Notification-Message-Id", payload.messageId)
	req.Header.Set("Notification-Message-Timestamp", payload.timestamp.Format(time.RFC3339Nano))
	req.Header.Set("Notification-Message-Type", notificationMessageTypeNotification)
	req.Header.Set("Content-Type", "application/json")

	if n.subscriberDetails.Secret != "" {
		// subscriber can validate the sender's identity by checking Notification-Message-Signature header in the request.
		// because it is hashed using the secret provided by them so that no one can create fake notifications.
		signature := n.createSignature(payload)
		req.Header.Set("Notification-Message-Signature", signature)
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("couldnt send the notification")
	}

	return nil
}

// createSignature creates the payload signature hash from messageId, timestamp and the payload
// by using a HMAC algorithm (SHA-256) with the secret provided from the subscriber.
func (n *notificationWebhookClient) createSignature(payload notificationPayload) string {
	h := hmac.New(sha256.New, []byte(n.subscriberDetails.Secret))

	buf, err := json.Marshal(payload)
	if err != nil {
		return "xx"
	}

	h.Write([]byte(fmt.Sprintf("%v%v%v", payload.messageId, payload.timestamp.Unix(), string(buf))))
	signature := fmt.Sprintf("sha256=%v", hex.EncodeToString(h.Sum(nil)))

	return signature
}
