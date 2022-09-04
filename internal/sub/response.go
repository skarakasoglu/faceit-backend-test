package sub

import "time"

// SubscribeResponse subscribe endpoint response model which containing the result of the subscriptions
// @Description  subscribe endpoint response model which containing the result of the subscriptions
type SubscribeResponse struct {
	Id        string    `json:"id"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
