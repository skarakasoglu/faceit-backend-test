package sub

import "time"

type SubscribeResponse struct {
	Id        string    `json:"id"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}
