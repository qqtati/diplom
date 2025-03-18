package common

type Response struct {
	Status      string      `json:"status"`
	Description string      `json:"description,omitempty"`
	Result      interface{} `json:"result,omitempty"`
}
