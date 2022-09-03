package sub

// SubscribeRequest subscribe endpoint request model contains the subscription parameters
// @Description subscribe endpoint request model contains the subscription parameters
type SubscribeRequest struct {
	Type     string `json:"type" binding:"required"`
	Callback string `json:"callback" binding:"required"`
	Secret   string `json:"secret"`
}
