package sub

type SubscribeRequest struct {
	Type     string `json:"type" binding:"required"`
	Callback string `json:"callback" binding:"required"`
	Secret   string `json:"secret"`
}
