package bitmexapi

type WsRequest struct {
	Operation string   `json:"op"`
	Args      []string `json:"args"`
}
