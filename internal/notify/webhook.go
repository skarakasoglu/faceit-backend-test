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

func verifySubscription(subscriber NotificationSubscriber) error {
	payload := verificationPayload{
		Id:        subscriber.Id,
		Challenge: uuid.New().String(),
		Callback:  subscriber.CallbackUrl,
		CreatedAt: subscriber.CreatedAt,
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

	req.Header.Set("Notification-Message-Type", notificationMessageTypeVerification)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		challenge := string(body)
		if challenge != payload.Challenge {
			return fmt.Errorf("invalid challenge provided: %v, expected: %v", challenge, payload.Challenge)
		}
	}

	return nil
}

func sendNotification(subscriber NotificationSubscriber, payload notificationPayload) error {
	buf, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	reqBody := bytes.NewReader(buf)
	req, err := http.NewRequest(http.MethodPost, subscriber.CallbackUrl, reqBody)
	if err != nil {
		return err
	}

	payload.timestamp = time.Now()

	req.Header.Set("Notification-Message-Id", payload.messageId)
	req.Header.Set("Notification-Message-Timestamp", fmt.Sprintf("%v", payload.timestamp.Unix()))
	req.Header.Set("Notification-Message-Type", notificationMessageTypeNotification)
	req.Header.Set("Content-Type", "application/json")

	if subscriber.Secret != "" {
		signature := createSignature(payload, subscriber.Secret)
		req.Header.Set("Notification-Message-Signature", signature)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("couldnt send the notification")
	}

	return nil
}

func createSignature(payload notificationPayload, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))

	buf, err := json.Marshal(payload)
	if err != nil {
		return "xx"
	}

	h.Write([]byte(fmt.Sprintf("%v%v%v", payload.messageId, payload.timestamp.Unix(), string(buf))))
	signature := fmt.Sprintf("sha256=%v", hex.EncodeToString(h.Sum(nil)))

	return signature
}
